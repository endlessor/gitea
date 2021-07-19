package lq

import (
	"errors"
	"fmt"
	"time"
)

func doProcess(action DomainActionInput) domainActionResult {
	//fmt.Println("run action", action)
	fmt.Println("run action <before sleep>")
	time.Sleep(3 * time.Second)
	fmt.Println("run action <after sleep>")
	str := "string"
	return domainActionResult{
		Stdout: &str,
		Error:  errors.New("hi"),
	}
}
