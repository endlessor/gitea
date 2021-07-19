package repo

import (
	"fmt"

	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/services/lq"
)

// lqBranchOp does the actual handling
func lqBranchOp(ctx *context.APIContext, currentUser string, httpMethod string, operation string, queryString string) {
	// TODO: validate inputs
	unsafeOwner := ctx.Params(":owner")
	// has permission?
	unsafeRepo := ctx.Params(":repo")
	unsafeBranch := ctx.Params(":branch")
	unsafeAction := ctx.Params(":action")
	payloadLocation := "TODO"
	s, err := lq.ProcessForResult(lq.DomainActionInput{
		DomainOwner:     unsafeOwner,
		TriggeredBy:     currentUser,
		DomainName:      unsafeRepo,
		Branch:          unsafeBranch,
		Action:          unsafeAction,
		HTTPMethod:      httpMethod,
		QueryString:     queryString,
		PayloadLocation: payloadLocation,
	})
	fmt.Println("Hi Context:", ctx, unsafeBranch, s, err)
	fmt.Println("Hi", ctx.Repo.BranchName, ctx.Repo.Repository.DefaultBranch, operation)
	// TODO: respond
	//ctx.SetLinkHeader(int(totalNumOfBranches), listOptions.PageSize)
	//ctx.Header().Set("X-Total-Count", fmt.Sprintf("%d", totalNumOfBranches))
	//ctx.Header().Set("Access-Control-Expose-Headers", "X-Total-Count, Link")
	//ctx.JSON(http.StatusOK, &apiBranches)
}

// LqBranchActionGet Do LQ operation for branch
func LqBranchActionGet(ctx *context.APIContext) {
	// swagger:operation GET /lq/{owner}/{repo}:{branch}/{action} lq lqBranchAction
	// ---
	// summary: Runs LQ action for specific branch
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: branch
	//   in: path
	//   description: branch name
	//   type: string
	//   required: true
	// - name: action
	//   in: path
	//   description: LQ action
	//   type: string
	//   required: true
	// responses:
	//   "200":

	// TODO: specify response type for swagger doc
	//     "$ref": "#/responses/BranchList"

	lqBranchOp(ctx, "currentUser", "get", "list", "queryString")
}

// LqDefaultBranchActionGet Do LQ operation for branch
func LqDefaultBranchActionGet(ctx *context.APIContext) {
	// swagger:operation GET /lq/{owner}/{repo}/{action} lq lqDefaultBranchAction
	// ---
	// summary: Runs LQ action for default branch
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: action
	//   in: path
	//   description: LQ action
	//   type: string
	//   required: true
	// responses:
	//   "200":

	// TODO: specify response type for swagger doc
	//     "$ref": "#/responses/BranchList"

	lqBranchOp(ctx, "currentUser", "get", "list", "queryString")
}
