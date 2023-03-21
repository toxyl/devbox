package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Commands Commands `mapstructure:"commands"`
	Limits   Limits   `mapstructure:"limits"`
	Options  Options  `mapstructure:"options"`
	Binds    Binds    `mapstructure:"binds"`
	Env      Env      `mapstructure:"env"`
}

func (c *Config) GetEnv() []string {
	return c.Env.Get()
}

func (c *Config) Save(file string) error {
	yamlConfig, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, yamlConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Load(file string) error {
	conf, err := DevboxFromFile(file)
	if err != nil {
		return err
	}
	c.Commands = conf.Commands
	c.Limits = conf.Limits
	c.Options = conf.Options
	c.Binds = conf.Binds
	c.Env = conf.Env
	return nil
}

func NewConfig() *Config {
	c := &Config{
		Commands: Commands{
			Start: "/usr/local/bin/start",
			Stop:  "/usr/local/bin/stop",
		},
		Limits: Limits{
			CPU: 0.1,
			Mem: LimitsMemory{
				Hard: 0.1,
				Soft: 0.75,
				Swap: 0.1,
			},
			PIDs: 1024,
		},
		Options: Options{
			MapUsersAndGroups: true,
			BindAll:           false,
		},
		Binds: Binds{},
		Env: Env{
			Shell:    "/bin/bash",
			Term:     "xterm",
			TermInfo: "/usr/share/terminfo/",
			Path:     "/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/snap/bin:/usr/sandbox/:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games",
			Home:     "/root",
			Vars:     EnvVars{},
		},
	}
	return c
}
