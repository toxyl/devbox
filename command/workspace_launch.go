package command

import (
	"fmt"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/devip"
	"github.com/toxyl/devbox/tmux"
	"github.com/toxyl/glog"
)

func WorkspaceLaunch(arg ...string) error {
	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file := filepath.Join(path, ".workspace.yaml")
	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}

	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	commands := make([]string, len(w.Devboxes))
	for i, devbox := range w.Devboxes {
		dboxPath := getDevboxConfigPath(devbox.Name)
		if !fileExists(dboxPath) {
			forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(dboxPath)), core.EXIT_DEVBOX_NOT_FOUND)
		}
		must(config.DevboxToFile(&devbox.Config, dboxPath))
		arg := ""
		if devbox.Delay > 0 {
			arg = "sleep " + fmt.Sprint(devbox.Delay) + " ; "
		}
		commands[i] = arg + core.APP_NAME + " " + ENTER + " " + devbox.Name
	}

	for _, ip := range w.IPs {
		devip.Add(ip)
	}

	defer func() {
		for _, ip := range w.IPs {
			devip.Remove(ip)
		}
	}()

	err = tmux.SpawnWorkspace(true, commands...)
	// store current container config back into workspace config
	for i, devbox := range w.Devboxes {
		dboxPath := getDevboxConfigPath(devbox.Name)
		if !fileExists(dboxPath) {
			forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(dboxPath)), core.EXIT_DEVBOX_NOT_FOUND)
		}
		if err := devbox.Config.Load(dboxPath); err != nil {
			return err
		}
		w.Devboxes[i] = devbox
	}

	w.Save(file)

	return err
}