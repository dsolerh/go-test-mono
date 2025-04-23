package main

import (
	"fmt"
	"workpublish"
)

func main() {
	c, err := workpublish.LoadPublishConfig("publish.yml")
	fmt.Printf("err: %v\n", err)
	fmt.Printf("c: %#v\n", c)

	err = workpublish.UpdateWorkspacePackages()
	fmt.Printf("err: %v\n", err)
}
