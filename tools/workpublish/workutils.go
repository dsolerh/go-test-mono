package workpublish

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/samber/lo"
)

func UpdateWorkspacePackages(newPackages, oldPackages []string) error {
	useNew := lo.Map(newPackages, func(p string, _ int) string { return fmt.Sprintf("-use=%s", p) })
	dropOld := lo.Map(oldPackages, func(p string, _ int) string { return fmt.Sprintf("-dropuse=%s", p) })
	baseArgs := []string{"work", "edit"}
	fullArgs := append(baseArgs, dropOld...)
	fullArgs = append(fullArgs, useNew...)
	output, err := exec.Command("go", append(baseArgs, dropOld...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	return nil
}

func UpdatePackageMods(c *PublishConfig, packages []string) error {
	for _, pkgName := range packages {
		oldPkgName := c.Packages[pkgName].Path[2:]
		// rename package name
		data, err := os.ReadFile(pkgName + "/go.mod")
		if err != nil {
			return fmt.Errorf("error reading go.mod file: %w", err)
		}
		data = bytes.Replace(data, []byte(oldPkgName), []byte(pkgName), 1)
		err = os.WriteFile(pkgName+"/go.mod", data, 0)
		if err != nil {
			return fmt.Errorf("error writing go.mod file: %w", err)
		}
	}
	return nil
}
