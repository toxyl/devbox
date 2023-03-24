package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func WorkspaceStore(arg ...string) error {
	name := arg[0]
	path := getWorkspacePath(name)
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
	file := getWorkspaceConfigPath(name)
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
	err = tar.FromDir(path, tarfile)
	if err != nil {
		log.Error("Could not store workspace to %s: %s", glog.File(tarfile), glog.Error(err))
		return nil
	}

	// remove the image files to save diskspace
	for _, c := range w.Devboxes {
		err = os.Remove(c.Image)
		if err != nil {
			return err
		}
	}

	log.Success("Stored workspace to %s", glog.File(tarfile))
	return nil
}
