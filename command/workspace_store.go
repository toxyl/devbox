package command

import (
	"fmt"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func WorkspaceStore(arg ...string) error {
	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !fileExists(path) {
		return fmt.Errorf("the source directory %s does not exist", path)
	}

	tarfile := arg[1]
	tarfile, err = filepath.Abs(tarfile)
	if err != nil {
		return err
	}
	log.Success("Storing workspace to %s", glog.File(tarfile))
	file := filepath.Join(path, ".workspace.yaml")
	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	for _, c := range w.Devboxes {
		name := c.Name
		err = Store(name, c.Image)
		if err != nil {
			return err
		}
	}
	tar.FromDir(path, tarfile)
	log.Success("Stored workspace to %s", glog.File(tarfile))
	return nil
}
