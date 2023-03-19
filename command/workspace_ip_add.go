package command

import (
	"fmt"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/glog"
)

func WorkspaceIPAdd(arg ...string) error {
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

	for _, ip := range arg[1:] {
		w.AddIP(ip)
	}

	w.Save(file)

	log.Success("IPs added to %s", glog.File(file))
	return nil
}
