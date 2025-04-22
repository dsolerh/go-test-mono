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
	for _, use := range workFile.Use {
		fmt.Printf("use: %#v\n", use)
	}
}
