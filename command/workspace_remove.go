package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/glog"
)

func WorkspaceRemove(arg ...string) error {
	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file := filepath.Join(path, ".workspace.yaml")
	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}

	devbooxes := arg[1:]
	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	for _, d := range devbooxes {
		_ = os.RemoveAll(filepath.Join(path, d+".tar.gz"))
	}
	err = w.Save(file)
	if err != nil {
		return err
	}
	log.Success("Devboxes removed from %s", glog.File(file))
	return nil
}
