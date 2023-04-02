package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/gutils"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	StoragePath   string `mapstructure:"storage_path" yaml:"storage_path"`
	AdminUser     string `mapstructure:"admin_user" yaml:"admin_user"`
	AdminPassword string `mapstructure:"admin_password" yaml:"admin_password"`
}

func (c *AppConfig) Path() string {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(cdir, ".devbox.yaml")
}

func (c *AppConfig) TestCredentials(user, password string) bool {
	return user == c.AdminUser && gutils.StringToSha256(c.AdminPassword) == password
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
	c.AdminUser = conf.AdminUser
	c.AdminPassword = conf.AdminPassword
	return nil
}

func NewAppConfig() *AppConfig {
	ac := &AppConfig{
		StoragePath:   os.TempDir(),
		AdminUser:     "",
		AdminPassword: "",
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
