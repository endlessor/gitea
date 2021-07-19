package lq

import (
	"errors"
	"sync"
)

var (
	initializeOnce sync.Once
	stopOnce       sync.Once
	inputCh        chan inputWithOutputCh
	doneCh         chan bool
)

type DomainActionInput struct {
	DomainOwner     string
	TriggeredBy     string
	DomainName      string
	Branch          string
	Action          string
	HTTPMethod      string
	QueryString     string
	PayloadLocation string
}

func (a DomainActionInput) toDomainBranchName() string {
	// "/" is an illegal character for domain and branch naming
	return a.DomainOwner + "/" + a.DomainName + "/" + a.Branch
}

type inputWithOutputCh struct {
	Input    DomainActionInput
	OutputCh *chan domainActionResult
}

type domainActionResult struct {
	Stdout *string
	Error  error
}

type resultWithInput struct {
	inputWithOutputCh
	Result domainActionResult
}

// InitRunner spawns a worker that processes domain actions
func InitRunner() {
	initializeOnce.Do(func() {
		inputCh, doneCh = spawnRunner()
	})
}

// StopRunner terminates the runner.
// Not recoverable, will cause later inputs to panic, returns err for all remaining tasks.
func StopRunner() {
	stopOnce.Do(func() {
		close(doneCh)
	})
}

func submitPanicky(action DomainActionInput, recoveredErrCh chan error, maybeOutputCh *chan domainActionResult) {
	defer func() {
		if err := recover(); err != nil {
			recoveredErrCh <- err.(error)
		}
	}()
	inputCh <- inputWithOutputCh{
		Input:    action,
		OutputCh: maybeOutputCh,
	}
}

func Process(action DomainActionInput) error {
	submitErrCh := make(chan error)
	go submitPanicky(action, submitErrCh, nil)
	select {
	case err := <-submitErrCh:
		return err
	default:
		return nil
	}
}

func ProcessForResult(action DomainActionInput) (*string, error) {
	outputCh := make(chan domainActionResult)
	submitErrCh := make(chan error)
	go submitPanicky(action, submitErrCh, &outputCh)
	select {
	case err := <-submitErrCh:
		return nil, err
	case result := <-outputCh:
		return result.Stdout, result.Error
	}
}

func respondToOutputCh(result domainActionResult, maybeOutputCh *chan domainActionResult) {
	if maybeOutputCh != nil {
		outputCh := *maybeOutputCh
		outputCh <- result
		close(outputCh)
	}
}

func dumpInputs(inputCh chan inputWithOutputCh) {
	close(inputCh)
	for input := range inputCh {
		go respondToOutputCh(domainActionResult{
			Error: errors.New("system is shutting down"),
		}, input.OutputCh)
	}
}

func handleResult(lockedRepoBranches *map[string]bool, result *resultWithInput) {
	repoBranchName := result.Input.toDomainBranchName()
	go respondToOutputCh(result.Result, result.OutputCh)
	delete(*lockedRepoBranches, repoBranchName)
}

func handleResultsBlocking(resultCh chan *resultWithInput, lockedRepoBranches *map[string]bool) {
	for result := range resultCh {
		handleResult(lockedRepoBranches, result)
		if len(*lockedRepoBranches) == 0 {
			return
		}
	}
}

func spawnRunner() (chan inputWithOutputCh, chan bool) {
	inputCh := make(chan inputWithOutputCh)
	doneCh := make(chan bool)
	resultCh := make(chan *resultWithInput)
	lockedRepoBranches := map[string]bool{}
	go func() {
		for {
			select {
			case action := <-inputCh:
				repoBranchName := action.Input.toDomainBranchName()
				if lockedRepoBranches[repoBranchName] {
					go respondToOutputCh(domainActionResult{
						Error: errors.New("this domain is currently being used by other action"),
					}, action.OutputCh)
					continue
				}
				lockedRepoBranches[repoBranchName] = true
				go doProcessIntoCh(resultCh, action)
			case result := <-resultCh:
				handleResult(&lockedRepoBranches, result)
			case <-doneCh:
				go handleResultsBlocking(resultCh, &lockedRepoBranches)
				dumpInputs(inputCh)
				return
			}
		}
	}()
	return inputCh, doneCh
}

func doProcessIntoCh(resultCh chan *resultWithInput, action inputWithOutputCh) {
	resultCh <- &resultWithInput{
		inputWithOutputCh: action,
		Result:            doProcess(action.Input),
	}
}
