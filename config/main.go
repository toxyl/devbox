package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/toxyl/glog"
)

func fileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}

func parseWithViper(file string, config interface{}) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	if !fileExists(path) {
		return fmt.Errorf("there is no config at %s", glog.File(path))
	}
	cPath := filepath.Dir(path)
	file = filepath.Base(path)
	configNameElements := strings.Split(file, ".")
	file = strings.Join(configNameElements[:len(configNameElements)-1], ".")
	cType := filepath.Ext(path)
	if len(cType) > 0 {
		cType = cType[1:]
	}
	v := viper.New()
	v.SetConfigName(file)
	v.SetConfigType(cType)
	v.AddConfigPath(cPath)
	err = v.ReadInConfig()
	if err != nil {
		return err
	}
	if err := v.Unmarshal(config); err != nil {
		return err
	}
	return nil
}

func DevboxFromFile(file string) (Config, error) {
	c := NewConfig()
	err := parseWithViper(file, c)
	return *c, err
}

func DevboxToFile(c *Config, file string) error {
	return c.Save(file)
}
