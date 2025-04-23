package workpublish

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PublishConfig struct {
	configPath string
	Repo       string                  `yaml:"repo"`
	Packages   map[string]PackagePInfo `yaml:"packages"`
}

type PackagePInfo struct {
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
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

func (c *PublishConfig) UpdatePackagesVersion(updater func(string) string) {
	for _, pkg := range c.Packages {
		pkg.Version = updater(pkg.Version)
	}
}

func (c *PublishConfig) SaveConfig() error {
	f, err := os.Create(c.configPath)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(f).Encode(c)
}
