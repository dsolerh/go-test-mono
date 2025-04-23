package main

import (
	"fmt"
	"os/exec"
)

var message = `ci: update packages

+ package1
+ package2
+ package3
+ package4
`

func main() {
	output, err := exec.Command(
		"git",
		"commit",
		"-m",
		message,
	).CombinedOutput()
	fmt.Printf("output: %s\n", output)
	fmt.Printf("err: %v\n", err)
}
