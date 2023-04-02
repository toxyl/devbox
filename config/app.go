package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/gutils"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	StoragePath       string `mapstructure:"storage_path" yaml:"storage_path"`
	RepoAdminUser     string `mapstructure:"repo_admin_user" yaml:"repo_admin_user"`
	RepoAdminPassword string `mapstructure:"repo_admin_password" yaml:"repo_admin_password"`
	RepoAddress       string `mapstructure:"repo_address" yaml:"repo_address"`
	RepoPath          string `mapstructure:"repo_path" yaml:"repo_path"`
}

func (c *AppConfig) Path() string {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(cdir, ".devbox.yaml")
}

func (c *AppConfig) TestCredentials(user, password string) bool {
	return user == c.RepoAdminUser && gutils.StringToSha256(c.RepoAdminPassword) == password
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
	c.RepoAdminUser = conf.RepoAdminUser
	c.RepoAdminPassword = conf.RepoAdminPassword
	c.RepoAddress = conf.RepoAddress
	c.RepoPath = conf.RepoPath
	return nil
}

func NewAppConfig() *AppConfig {
	ac := &AppConfig{
		StoragePath:       os.TempDir(),
		RepoAdminUser:     "admin",
		RepoAdminPassword: "admin",
		RepoAddress:       "127.0.0.1:438",
		RepoPath:          filepath.Join(os.TempDir(), "repo"),
	}
	return ac
}

func OpenAppConfig() (*AppConfig, error) {
	// check if AppConfig exists, else create default one
	ac := NewAppConfig()
	_, err := os.Open(ac.Path())
	if err != nil {
		if err := ac.Save(); err != nil {
			return nil, fmt.Errorf("could not create default app config")
		}

	}
	err = ac.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load app config")
	}
	return ac, nil
}
