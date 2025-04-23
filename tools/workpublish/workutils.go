package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func AddPackagesToWorspace(pmap PackagesMap, packagesName []string) error {
	// exec go work use <name-of-the-package>
	baseArgs := []string{"work", "use"}
	cmd := exec.Command("go", append(baseArgs, packagesName...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to workspace: %w, output: %s\n", err, output)
	}

	dropOld := make([]string, 0, 2*len(pmap))
	for _, pkgInfo := range pmap {
		dropOld = append(dropOld, fmt.Sprintf("-dropreplace=%s", pkgInfo.CurrentPath))
		dropOld = append(dropOld, fmt.Sprintf("-dropuse=%s", pkgInfo.CurrentPath))
	}
	baseArgs = []string{"work", "edit"}
	output, err = exec.Command("go", append(baseArgs, dropOld...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	for _, packageName := range packagesName {
		oldPackageName := pmap[packageName].CurrentPath
		// rename package name
		data, err := os.ReadFile(packageName + "/go.mod")
		if err != nil {
			return fmt.Errorf("error reading go.mod file: %w", err)
		}
		data = bytes.Replace(data, []byte(oldPackageName), []byte(packageName), 1)
		err = os.WriteFile(packageName+"/go.mod", data, 0)
		if err != nil {
			return fmt.Errorf("error writing go.mod file: %w", err)
		}
	}
	return nil
}

func RemovePackagesFromWorspace(pmap PackagesMap, packagesName []string) error {
	updateUses := make([]string, 0, len(packagesName)+len(pmap))
	for _, pkgName := range packagesName {
		updateUses = append(updateUses, fmt.Sprintf("-dropuse=%s", pkgName))
	}
	for _, pkgInfo := range pmap {
		updateUses = append(updateUses, fmt.Sprintf("-use=%s", pkgInfo.CurrentPath))
	}

	baseArgs := []string{"work", "edit"}
	output, err := exec.Command("go", append(baseArgs, updateUses...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	return nil
}

func UpdatePackagesVersions(pmap PackagesMap) error {
	updateReplaceVersion := make([]string, 0, len(pmap))
	for _, pkgInfo := range pmap {
		oldPath := pkgInfo.CurrentPath
		newPath := pkgInfo.PublishPath
		version := pkgInfo.Version
		updateReplaceVersion = append(updateReplaceVersion, fmt.Sprintf("-replace=%s=%s@%s", oldPath, newPath, version))
	}
	baseArgs := []string{"work", "edit"}
	output, err := exec.Command("go", append(baseArgs, updateReplaceVersion...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error updating packages version in workspace: %w, output: %s\n", err, output)
	}
	return nil
}
