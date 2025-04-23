package main

import (
	"flag"
	"log"
	"workpublish"
)

// Example usage
func main() {
	// check if there're uncommitted changes in the repo
	uncomm, err := workpublish.HasUncommittedChanges()
	if err != nil {
		log.Fatal(err)
	}
	if uncomm {
		log.Fatal("uncommitted changes in the repo")
	}

	v := flag.String("v", "patch", "the update to the package version (mayor|minor|patch)")
	fname := flag.String("f", "publish.yml", "the config file for the update")
	push := flag.Bool("push", false, "if set to true will push the changes to the git repo")
	flag.Parse()

	var vupdater = workpublish.SemverUpdater(*v)

	config, err := workpublish.LoadPublishConfig(*fname)
	if err != nil {
		log.Fatal(err)
	}

	pkgNames := flag.Args()
	if len(pkgNames) == 0 {
		pkgNames = config.AllPackageNames()
	}
	log.Printf("current packages: %v\n", config.GetTagVersions(config.AllPackageNames()))
	// update versions
	config.UpdatePackagesVersion(pkgNames, vupdater)
	tagVersions := config.GetTagVersions(pkgNames)
	oldPackages := config.GetOldPackages()

	log.Printf("publishing packages: %v\n", tagVersions)
	log.Println("copiying packages to root...")
	// copy the packages to the root of the workspace
	if err := workpublish.CopyPackagesToRoot(config, pkgNames); err != nil {
		log.Fatal(err)
	}

	log.Println("updating workspace packages...")
	// remove the old packages from the go.work and add the new
	if err := workpublish.UpdateWorkspacePackages(pkgNames, oldPackages); err != nil {
		log.Println("removing packages from root...")
		// remove the packages from root (cleanup)
		if err2 := workpublish.RemovePackagesFromRoot(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	// update the go.mod files of the new packages
	if err := workpublish.UpdatePackageMods(config, pkgNames); err != nil {
		log.Println("removing packages from root...")
		// remove the packages from root (cleanup)
		if err2 := workpublish.RemovePackagesFromRoot(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Println("reverting workspace packages...")
		// revert the workspace packages
		if err3 := workpublish.UpdateWorkspacePackages(oldPackages, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	// commit all changes
	if err := workpublish.CommitChanges(workpublish.GetPublishCommitMessage(pkgNames)); err != nil {
		log.Println("removing packages from root...")
		// remove the packages from root (cleanup)
		if err2 := workpublish.RemovePackagesFromRoot(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Println("reverting workspace packages...")
		// revert the workspace packages
		if err3 := workpublish.UpdateWorkspacePackages(oldPackages, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	// tag the versions
	if err := workpublish.TagPackagesVersion(tagVersions); err != nil {
		log.Println("removing packages from root...")
		// remove the packages from root (cleanup)
		if err2 := workpublish.RemovePackagesFromRoot(pkgNames); err2 != nil {
			log.Println(err2)
		}
		log.Println("reverting workspace packages...")
		// revert the workspace packages
		if err3 := workpublish.UpdateWorkspacePackages(oldPackages, pkgNames); err3 != nil {
			log.Println(err3)
		}
		// TODO: remove commits
		log.Fatal(err)
	}

	log.Println("removing packages from root...")
	if err := workpublish.RemovePackagesFromRoot(pkgNames); err != nil {
		log.Println("reverting workspace packages...")
		// revert the workspace packages
		if err3 := workpublish.UpdateWorkspacePackages(oldPackages, pkgNames); err3 != nil {
			log.Println(err3)
		}
		log.Fatal(err)
	}

	log.Println("reverting workspace packages...")
	// revert the workspace packages
	if err := workpublish.UpdateWorkspacePackages(oldPackages, pkgNames); err != nil {
		log.Fatal(err)
	}

	// save the version changes
	if err = config.SaveConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println("committing reverted changes...")
	if err := workpublish.CommitChanges(workpublish.CleanupCommit); err != nil {
		log.Fatal(err)
	}

	if *push {
		log.Println("pushing changes...")
		if err := workpublish.PushChanges(tagVersions); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("publish completed...")
}
