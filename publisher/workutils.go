package publisher

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

	dropReplaces := make([]string, 0, len(packagesName))
	for _, pkgName := range packagesName {
		oldPath := pmap[pkgName].OldPath
		dropReplaces = append(dropReplaces, fmt.Sprintf("-dropreplace=%s", oldPath))
	}
	baseArgs = []string{"work", "edit"}
	output, err = exec.Command("go", append(baseArgs, dropReplaces...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	for _, packageName := range packagesName {
		oldPackageName := pmap[packageName].OldPath
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

func RemovePackagesFromWorspace(packagesName []string) error {
	dropUses := make([]string, 0, len(packagesName))
	for _, pkgName := range packagesName {
		dropUses = append(dropUses, fmt.Sprintf("-dropuse=%s", pkgName))
	}
	baseArgs := []string{"work", "edit"}
	output, err := exec.Command("go", append(baseArgs, dropUses...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	return nil
}

func UpdatePackagesVersions(pmap PackagesMap, packagesName []string) error {
	updateReplaceVersion := make([]string, 0, len(packagesName))
	for _, pkgName := range packagesName {
		oldPath := pmap[pkgName].OldPath
		newPath := pmap[pkgName].NewPath
		version := pmap[pkgName].Version
		updateReplaceVersion = append(updateReplaceVersion, fmt.Sprintf("-replace=%s=%s@v%s", oldPath, newPath, version))
	}
	baseArgs := []string{"work", "edit"}
	output, err := exec.Command("go", append(baseArgs, updateReplaceVersion...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error updating packages version in workspace: %w, output: %s\n", err, output)
	}
	return nil
}
