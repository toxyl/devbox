package command

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func WorkspaceAdd(arg ...string) error {
	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file := filepath.Join(path, ".workspace.yaml")
	devboxes := arg[1:]
	configs := []config.Config{}
	delays := []int{}
	names := []string{}
	for _, dbox := range devboxes {
		elems := strings.Split(dbox, ":")
		dboxName := elems[0]
		dboxPath := getDevboxConfigPath(dboxName)
		if !fileExists(dboxPath) {
			forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(dboxPath)), core.EXIT_DEVBOX_NOT_FOUND)
		}
		dboxConfig, err := config.DevboxFromFile(dboxPath)
		delay := 0
		if len(elems) == 2 {
			delay, err = strconv.Atoi(elems[1])
			if err != nil {
				return errors.Join(err, fmt.Errorf("delay argument must be an int"))
			}
		}
		if err != nil {
			return err
		}
		configs = append(configs, *dboxConfig)
		delays = append(delays, delay)
		names = append(names, dboxName)
	}
	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	for i, c := range configs {
		name := names[i]
		image := filepath.Join(path, name+".tar.gz")
		err = Store(name, image)
		if err != nil {
			return err
		}
		delay := delays[i]
		w.Add(name, image, int64(delay), c)
	}
	err = w.Save(file)
	if err != nil {
		return err
	}
	log.Success("Devboxes added to %s", glog.File(file))
	return nil
}
