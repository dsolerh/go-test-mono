package workpublish

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type PublishConfig struct {
	configPath string
	Repo       string                  `yaml:"repo"`
	Root       string                  `yaml:"root"`
	Packages   map[string]*PackageInfo `yaml:"packages"`
}

type PackageInfo struct {
	WorkName string `yaml:"work_name"`
	PkgName  string `yaml:"pkg_name"`
	Version  string `yaml:"version"`
}

func LoadPublishConfig(fname string) (*PublishConfig, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	config := new(PublishConfig)
	config.configPath = fname

	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *PublishConfig) AllPackageNames() []string {
	packageNames := make([]string, 0, len(c.Packages))
	for pkgName := range c.Packages {
		packageNames = append(packageNames, pkgName)
	}
	return packageNames
}

func (c *PublishConfig) UpdatePackagesVersion(packages []string, updater func(string) string) {
	for _, pkgName := range packages {
		pkg := c.Packages[pkgName]
		pkg.Version = updater(pkg.Version)
	}
}

func (c *PublishConfig) GetTagVersions(packages []string) []string {
	tags := make([]string, 0, len(packages))
	for _, pkgName := range packages {
		tags = append(tags, fmt.Sprintf("%s/%s", pkgName, c.Packages[pkgName].Version))
	}
	return tags
}

func (c *PublishConfig) GetOldPackages() []string {
	packages := make([]string, 0, len(c.Packages))
	for _, pkg := range c.Packages {
		packages = append(packages, pkg.WorkName)
	}
	return packages
}

func (c *PublishConfig) SaveConfig() error {
	f, err := os.Create(c.configPath)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(f).Encode(c)
}
