package main

import (
	"flag"
	"log"
	"workpublish"
)

// Example usage
func main() {
	v := flag.String("v", "patch", "the update to the package version (mayor|minor|patch)")
	flag.Parse()

	var vupdater = workpublish.SemverUpdater(*v)

	workFile, err := workpublish.ParseGoWorkFile("go.work")
	if err != nil {
		panic(err)
	}

	pmap := workpublish.MakePackagesMap(workFile)

	pkgNames := flag.Args()
	if len(pkgNames) == 0 {
		for pkey := range pmap {
			pkgNames = append(pkgNames, pkey)
		}
	}

	if err := workpublish.CopyDirectories(pmap, pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := workpublish.AddPackagesToWorspace(pmap, pkgNames); err != nil {
		if err2 := workpublish.RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	if err := workpublish.CommitAndTagChanges(pmap, pkgNames, vupdater); err != nil {
		if err2 := workpublish.RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		if err3 := workpublish.RemovePackagesFromWorspace(pmap, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	if err := workpublish.UpdatePackagesVersions(pmap); err != nil {
		log.Fatal(err)
	}

	if err := workpublish.RemoveAllPackages(pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := workpublish.RemovePackagesFromWorspace(pmap, pkgNames); err != nil {
		log.Println(err)
	}

	if err := workpublish.CleanUpCommit(pkgNames); err != nil {
		log.Println(err)
	}

	if err := workpublish.PushChanges(pmap, pkgNames); err != nil {
		log.Println(err)
	}
}
