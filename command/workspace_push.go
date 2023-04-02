package command

import (
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func WorkspacePush(arg ...string) error {
	name := arg[0]
	path := getWorkspacePath(name)
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !fileExists(path) {
		forceExit("The source directory does not exist: "+glog.File(path), core.EXIT_WORKSPACE_PUSH_FAILED)
	}
	tarname := name + ".tar.gz"
	tarfile := filepath.Join(getStorageDir(), tarname)
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
		forceExit("Could not store workspace to "+glog.File(tarfile)+": "+glog.Error(err), core.EXIT_WORKSPACE_PUSH_FAILED)
	}

	// remove the image files to save diskspace
	for _, c := range w.Devboxes {
		err = os.Remove(c.Image)
		if err != nil {
			return err
		}
	}

	log.Success("Stored workspace to %s", glog.File(tarfile))

	return RepoUpload(tarfile, "workspace_"+tarname)
}
