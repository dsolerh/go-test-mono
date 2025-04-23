package workpublish

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func UpdateWorkspacePackages() error {
	output, err := exec.Command("go", "work", "use", "-r", ".").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing packages from workspace: %w, output: %s\n", err, output)
	}
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
