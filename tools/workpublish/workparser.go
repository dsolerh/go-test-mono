package main

import (
	"os"
	"strings"

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

type PackageInfo struct {
	CurrentPath string // the path the package has in the current project structure
	PublishPath string // the package name and new path
	Version     string // the package current version
}

type PackagesMap map[string]*PackageInfo

func MakePackagesMap(work *modfile.WorkFile) PackagesMap {
	pmap := make(PackagesMap, len(work.Replace))
	for _, replace := range work.Replace {
		parts := strings.Split(replace.Old.Path, "/")
		pkey := parts[len(parts)-1]
		pmap[pkey] = &PackageInfo{
			CurrentPath: replace.New.Path,
			PublishPath: replace.Old.Path,
			Version:     replace.Old.Version,
		}
	}
	return pmap
}
