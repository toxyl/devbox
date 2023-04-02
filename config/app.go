package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type RepoUserConfig struct {
	Admin    bool   `mapstructure:"admin" yaml:"admin"`
	Name     string `mapstructure:"name" yaml:"name"`
	Password string `mapstructure:"password" yaml:"password"`
}

type RepoServerConfig struct {
	Path    string           `mapstructure:"path" yaml:"path"`
	Address string           `mapstructure:"address" yaml:"address"`
	Users   []RepoUserConfig `mapstructure:"users" yaml:"users"`
}

type RepoClientConfig struct {
	Address  string `mapstructure:"address" yaml:"address"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
}

type RepoConfig struct {
	Server RepoServerConfig `mapstructure:"server" yaml:"server"`
	Client RepoClientConfig `mapstructure:"client" yaml:"client"`
}

type AppConfig struct {
	StoragePath string     `mapstructure:"storage_path" yaml:"storage_path"`
	Repo        RepoConfig `mapstructure:"repo" yaml:"repo"`
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
	c.Repo.Client.Address = conf.Repo.Client.Address
	c.Repo.Client.User = conf.Repo.Client.User
	c.Repo.Client.Password = conf.Repo.Client.Password
	c.Repo.Server.Address = conf.Repo.Server.Address
	c.Repo.Server.Path = conf.Repo.Server.Path
	c.Repo.Server.Users = conf.Repo.Server.Users

	return nil
}

func NewAppConfig() *AppConfig {
	ac := &AppConfig{
		StoragePath: os.TempDir(),
		Repo: RepoConfig{
			Server: RepoServerConfig{
				Path:    filepath.Join(os.TempDir(), "repo"),
				Address: "127.0.0.1:438",
				Users: []RepoUserConfig{
					{
						Admin:    true,
						Name:     "admin",
						Password: fmt.Sprint(time.Now().Unix()),
					},
				},
			},
			Client: RepoClientConfig{
				Address:  "127.0.0.1:438",
				User:     "user",
				Password: fmt.Sprint(time.Now().Unix()),
			},
		},
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
