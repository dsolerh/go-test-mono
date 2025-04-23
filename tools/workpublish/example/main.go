package main

import (
	"fmt"
	"os/exec"
)

func main() {
	output, err := exec.Command(
		"go",
		"work",
		"edit",
		"-use=./packages/core/foo",
		"-use=./tools/workpublish",
	).CombinedOutput()
	fmt.Printf("output: %s\n", output)
	fmt.Printf("err: %v\n", err)
}
