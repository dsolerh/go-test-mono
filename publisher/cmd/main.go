package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dsolerh/go-test-mono/packages/publisher"
)

func CommitAndTagChanges(packageName, packageVersion string) error {
	// add the changes
	output, err := exec.Command("git", "add", "--", packageName, "go.work").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to git: %w, output: %s\n", err, output)
	}

	// commit the changes
	output, err = exec.Command("git", "commit", "-m", fmt.Sprintf(`"ci: publish package %s"`, packageName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error commiting the changes: %w, output: %s\n", err, output)
	}

	// tag the commit
	tag := fmt.Sprintf("%s/%s", packageName, packageVersion)
	output, err = exec.Command("git", "tag", "--", tag, "HEAD").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging the changes: %w, output: %s\n", err, output)
	}

	// push the changes
	output, err = exec.Command("git", "push").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the changes: %w, output: %s\n", err, output)
	}

	// push the tag
	output, err = exec.Command("git", "push", "origin", "tag", tag).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the tag: %v, output: %s\n", err, output)
	}

	return nil
}

func CleanUp(packageName string) error {
	// remove the package from the root
	if err := os.RemoveAll(packageName); err != nil {
		return fmt.Errorf("Error removing package directory: %w\n", err)
	}
	return nil
}

func RevertCommit(packageName string) error {
	output, err := exec.Command("go", "work", "use", "-r", ".").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to workspace: %w, output: %s\n", err, output)
	}

	// add the changes
	output, err = exec.Command("git", "add", "--", packageName, "go.work").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to git: %v, output: %s\n", err, output)
	}

	// commit the changes
	output, err = exec.Command("git", "commit", "-m", fmt.Sprintf(`"ci: cleanup publish of package %s"`, packageName)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error commiting the changes: %v, output: %s\n", err, output)
	}

	// push the changes
	output, err = exec.Command("git", "push").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the changes: %v, output: %s\n", err, output)
	}

	return nil
}

func main() {
	pathToPackage := flag.String("pp", "", "the path to the package to publish")
	nameOfPackage := flag.String("pn", "", "the name of the package")
	packageVersion := flag.String("pv", "", "the version of the package to publish")
	flag.Parse()

	if err := publisher.CopyDirectory(*pathToPackage, *nameOfPackage); err != nil {
		log.Fatal(err)
	}

	if err := publisher.AddPackagesToWorspace(nil, []string{*nameOfPackage}); err != nil {
		if err2 := CleanUp(*nameOfPackage); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	if err := CommitAndTagChanges(*nameOfPackage, *packageVersion); err != nil {
		if err2 := CleanUp(*nameOfPackage); err2 != nil {
			log.Println(err2)
		}
		log.Fatal(err)
	}

	if err := CleanUp(*nameOfPackage); err != nil {
		log.Println(err)
	}

	if err := RevertCommit(*nameOfPackage); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("package %s published with version %s\n", *nameOfPackage, *packageVersion)
}
