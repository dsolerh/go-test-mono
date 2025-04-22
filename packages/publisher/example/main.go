package main

import (
	"fmt"
	"workparser"
)

// Example usage
func main() {
	workFile, err := workparser.ParseGoWorkFile("go.work")
	if err != nil {
		panic(err)
	}

	// Access workspace information
	for _, replace := range workFile.Replace {
		fmt.Printf("replace: %#v\n", replace)
	}
}
