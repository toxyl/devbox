package command

import (
	"fmt"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/glog"
)

func WorkspaceIPAdd(arg ...string) error {
	name := arg[0]
	file := getWorkspaceConfigPath(name)

	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}

	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}

	for _, ip := range arg[1:] {
		w.AddIP(ip)
	}

	w.Save(file)

	log.Success("IPs added to %s", glog.File(file))
	return nil
}
