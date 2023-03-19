package config

import (
	"fmt"
	"strings"
)

type Env struct {
	Shell    string  `mapstructure:"shell"`
	Term     string  `mapstructure:"term"`
	TermInfo string  `mapstructure:"term_info" yaml:"term_info"`
	Path     string  `mapstructure:"path"`
	Home     string  `mapstructure:"home"`
	Vars     EnvVars `mapstructure:"vars"`
}

func (c *Env) Get() []string {
	env := []string{
		"SHELL=" + c.Shell,
		"TERM=" + c.Term,
		"TERMINFO=" + c.TermInfo,
		"PATH=" + c.Path,
		"HOME=" + c.Home,
	}
	for k, v := range c.Vars {
		env = append(env, strings.ToUpper(strings.TrimSpace(k))+"="+strings.TrimSpace(fmt.Sprint(v)))
	}
	return env
}
