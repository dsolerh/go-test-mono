package main

import (
	"flag"
	"log"
)

// Example usage
func main() {
	v := flag.String("v", "patch", "the update to the package version (mayor|minor|patch)")
	flag.Parse()

	var vupdater = SemverUpdater(*v)

	workFile, err := ParseGoWorkFile("go.work")
	if err != nil {
		panic(err)
	}

	pmap := MakePackagesMap(workFile)

	pkgNames := flag.Args()
	if len(pkgNames) == 0 {
		for pkey := range pmap {
			pkgNames = append(pkgNames, pkey)
		}
	}

	if err := CopyDirectories(pmap, pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := AddPackagesToWorspace(pmap, pkgNames); err != nil {
		if err2 := RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	if err := CommitAndTagChanges(pmap, pkgNames, vupdater); err != nil {
		if err2 := RemoveAllPackages(pkgNames); err2 != nil {
			log.Println(err2)
		}
		if err3 := RemovePackagesFromWorspace(pmap, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	if err := UpdatePackagesVersions(pmap); err != nil {
		log.Fatal(err)
	}

	if err := RemoveAllPackages(pkgNames); err != nil {
		log.Fatal(err)
	}

	if err := RemovePackagesFromWorspace(pmap, pkgNames); err != nil {
		log.Println(err)
	}

	if err := CleanUpCommit(pkgNames); err != nil {
		log.Println(err)
	}

	if err := PushChanges(pmap, pkgNames); err != nil {
		log.Println(err)
	}
}
