package workparser

import (
	"os"

	"golang.org/x/mod/modfile"
)

func ParseGoWorkFile(path string) (*modfile.WorkFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the go.work file
	return modfile.ParseWork(path, data, nil)
}
