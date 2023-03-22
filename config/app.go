package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	StoragePath string `mapstructure:"storage_path" yaml:"storage_path"`
}

func (c *AppConfig) Path() string {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(cdir, ".devbox.yaml")
}

func (c *AppConfig) Save() error {
	file := c.Path()
	yamlConfig, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, yamlConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) Load() error {
	conf, err := GlobalFromFile(c.Path())
	if err != nil {
		return err
	}

	c.StoragePath = conf.StoragePath
	return nil
}

func NewAppConfig() *AppConfig {
	ac := &AppConfig{
		StoragePath: os.TempDir(),
	}
	return ac
}
