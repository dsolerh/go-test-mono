package workpublish

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/samber/lo"
)

func UpdateWorkspacePackages(newPackages, oldPackages []string) error {
	useNew := lo.Map(newPackages, func(p string, _ int) string { return fmt.Sprintf("-use=%s", p) })
	dropOld := lo.Map(oldPackages, func(p string, _ int) string { return fmt.Sprintf("-dropuse=%s", p) })
	baseArgs := []string{"work", "edit"}
	fullArgs := append(baseArgs, dropOld...)
	fullArgs = append(fullArgs, useNew...)
	log.Printf("running: 'go %s'", strings.Join(fullArgs, " "))
	output, err := exec.Command("go", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
	log.Printf("output: '%s'", output)
	return nil
}

func UpdatePackageMods(c *PublishConfig, packages []string) error {
	for _, pkgn := range packages {
		workName := c.Packages[pkgn].WorkName
		pkgName := c.Packages[pkgn].PkgName
		// rename package name
		gomodPath := path.Join(pkgName, "go.mod")
		log.Printf("updating %s\n", gomodPath)
		data, err := os.ReadFile(gomodPath)
		if err != nil {
			return fmt.Errorf("error reading go.mod file: %w", err)
		}
		log.Printf("replace %s with %s", workName, pkgName)
		data = bytes.Replace(data, []byte(workName), []byte(pkgName), 1)
		err = os.WriteFile(gomodPath, data, 0)
		if err != nil {
			return fmt.Errorf("error writing go.mod file: %w", err)
		}
	}
	return nil
}
