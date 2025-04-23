package workpublish

import (
	"fmt"
	"os"
	"os/exec"
)

func copyDirectory(src, dst string) error {
	// The -r flag makes cp recursive
	// The -p flag preserves mode, ownership, and timestamps
	cmd := exec.Command("cp", "-rp", src, dst)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error copying directory: %v, output: %s", err, output)
	}

	return nil
}

func CopyDirectories(pmap PackagesMap, packages []string) error {
	for _, pkg := range packages {
		src := pmap[pkg].CurrentPath
		dst := pkg
		if err := copyDirectory(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func CopyPackagesToRoot(c *PublishConfig, packages []string) error {
	for _, pkg := range packages {
		src := c.Packages[pkg].Path
		dst := pkg
		if err := copyDirectory(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func RemovePackagesFromRoot(packagesName []string) error {
	for _, packageName := range packagesName {
		// remove the package from the root
		if err := os.RemoveAll(packageName); err != nil {
			return fmt.Errorf("Error removing package directory: %w\n", err)
		}
	}
	return nil
}
