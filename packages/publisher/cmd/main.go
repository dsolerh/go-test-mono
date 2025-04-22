package main

import (
	"flag"
	"log"

	"github.com/dsolerh/go-test-mono/packages/publisher"
)

// Example usage
func main() {
	v := flag.String("v", "patch", "the update to the package version (mayor|minor|patch)")
	flag.Parse()

	var vupdater = publisher.SemverUpdater(*v)

	workFile, err := publisher.ParseGoWorkFile("go.work")
	if err != nil {
		panic(err)
	}

	pmap := publisher.MakePackagesMap(workFile)

	pkgNames := flag.Args()
	if len(pkgNames) == 0 {
		for pkey := range pmap {
			pkgNames = append(pkgNames, pkey)
		}
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

	if err := publisher.CommitAndTagChanges(pmap, pkgNames, vupdater); err != nil {
		if err2 := publisher.RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		if err3 := publisher.RemovePackagesFromWorspace(pmap, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	if err := publisher.UpdatePackagesVersions(pmap); err != nil {
		log.Fatal(err)
	}

	if err := publisher.RemoveAllPackages(pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := publisher.RemovePackagesFromWorspace(pmap, pkgNames); err != nil {
		log.Println(err)
	}

	if err := publisher.CleanUpCommit(pkgNames); err != nil {
		log.Println(err)
	}

	if err := publisher.PushChanges(pmap, pkgNames); err != nil {
		log.Println(err)
	}
}
