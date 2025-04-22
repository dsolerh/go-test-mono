package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func CopyDirectory(src, dst string) error {
	// The -r flag makes cp recursive
	// The -p flag preserves mode, ownership, and timestamps
	cmd := exec.Command("cp", "-rp", src, dst)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error copying directory: %v, output: %s", err, output)
	}

	return nil
}

func RemoveDir() {}

func main() {
	pathToPackage := flag.String("pp", "", "the path to the package to publish")
	nameOfPackage := flag.String("pn", "", "the name of the package")
	packageVersion := flag.String("pv", "", "the version of the package to publish")
	flag.Parse()

	// copy the package to the root of the project
	err := CopyDirectory(*pathToPackage, *nameOfPackage)
	if err != nil {
		fmt.Printf("Error copying directory: %v\n", err)
		return
	}

	// exec go work use <name-of-the-package>
	cmd := exec.Command("go", "work", "use", *nameOfPackage)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error adding package to workspace: %v, output: %s\n", err, output)
		return
	}

	// add the changes
	cmd = exec.Command("git", "add", "--", *nameOfPackage, "go.work")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error adding package to git: %v, output: %s\n", err, output)
		return
	}

	// commit the changes
	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf(`"ci: publish package %s"`, *nameOfPackage))
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error commiting the changes: %v, output: %s\n", err, output)
		return
	}

	// tag the commit
	tag := fmt.Sprintf("%s/%s", *nameOfPackage, *packageVersion)
	cmd = exec.Command("git", "tag", "--", tag, "HEAD")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error tagging the changes: %v, output: %s\n", err, output)
		return
	}

	// push the changes
	cmd = exec.Command("git", "push")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error pushing the changes: %v, output: %s\n", err, output)
		return
	}

	// push the tag
	cmd = exec.Command("git", "push", "origin", "tag", tag)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error pushing the tag: %v, output: %s\n", err, output)
		return
	}

	// remove the package from the root
	err = os.RemoveAll(*nameOfPackage)
	if err != nil {
		fmt.Printf("Error removing package directory: %v\n", err)
		return
	}

	cmd = exec.Command("go", "work", "use", *nameOfPackage)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error adding package to workspace: %v, output: %s\n", err, output)
		return
	}

	// add the changes
	cmd = exec.Command("git", "add", "--", *nameOfPackage, "go.work")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error adding package to git: %v, output: %s\n", err, output)
		return
	}

	// commit the changes
	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf(`"ci: cleanup publish of package %s"`, *nameOfPackage))
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error commiting the changes: %v, output: %s\n", err, output)
		return
	}

	// push the changes
	cmd = exec.Command("git", "push")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error pushing the changes: %v, output: %s\n", err, output)
		return
	}

	fmt.Printf("package %s published with version %s\n", *nameOfPackage, *packageVersion)
}
