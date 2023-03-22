package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/glog"
)

func WorkspaceRemove(arg ...string) error {
	name := arg[0]
	file := getWorkspaceConfigPath(name)

	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}

	devboxes := arg[1:]
	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	path := getWorkspacePath(name)
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	for _, d := range devboxes {
		_ = os.RemoveAll(filepath.Join(path, d+".tar.gz"))
		w.Remove(d)
	}
	err = w.Save(file)
	if err != nil {
		return err
	}
	log.Success("Devboxes removed from %s", glog.File(file))
	return nil
}
