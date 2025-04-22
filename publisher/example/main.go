package main

import (
	"log"

	"github.com/dsolerh/go-test-mono/packages/publisher"
)

// Example usage
func main() {
	workFile, err := publisher.ParseGoWorkFile("go.work")
	if err != nil {
		panic(err)
	}

	pmap := publisher.MakePackagesMap(workFile)

	pkgNames := make([]string, 0, len(pmap))
	for pkey := range pmap {
		pkgNames = append(pkgNames, pkey)
	}

	if err := publisher.CopyDirectories(pmap, pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := publisher.AddPackagesToWorspace(pmap, pkgNames); err != nil {
		if err2 := publisher.RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	if err := publisher.CommitAndTagChanges(pmap, pkgNames); err != nil {
		if err2 := publisher.RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		if err3 := publisher.RemovePackagesFromWorspace(pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	if err := publisher.RemoveAllPackages(pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := publisher.RemovePackagesFromWorspace(pkgNames); err != nil {
		log.Println(err)
	}

	if err := publisher.CleanUpCommit(pkgNames); err != nil {
		log.Println(err)
	}

	if err := publisher.PushChanges(); err != nil {
		log.Println(err)
	}
}
