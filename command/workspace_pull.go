package command

import (
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func WorkspacePull(arg ...string) error {
	name := arg[0]
	dstDir := getWorkspacePath(name)
	dstDir, err := filepath.Abs(dstDir)
	if err != nil {
		return err
	}

	tarFileLocal := name + ".tar.gz"
	tarFile := "workspace_" + tarFileLocal
	isNew := RepoDownload(tarFile, tarFileLocal)
	if !isNew {
		return nil // nothing to do then
	}
	tarFile = filepath.Join(getStorageDir(), tarFileLocal)
	if !fileExists(tarFile) {
		forceExit("The source file does not exist: "+glog.File(tarFile), core.EXIT_WORKSPACE_PULL_FAILED)
	}
	log.Success("Restoring workspace from %s", glog.File(tarFile))

	err = tar.ToDir(tarFile, dstDir)
	if err != nil {
		return err
	}

	file := getWorkspaceConfigPath(name)
	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	for i, c := range w.Devboxes {
		name := c.Name
		c.Image = filepath.Join(dstDir, filepath.Base(c.Image))
		err = Make(name, c.Image)
		if err != nil {
			return err
		}
		w.Devboxes[i].Image = c.Image
	}
	w.Path = dstDir
	w.Save(file)

	// remove the image files to save diskspace
	for _, c := range w.Devboxes {
		err = os.Remove(c.Image)
		if err != nil {
			return err
		}
	}

	log.Success("Restored workspace to %s", glog.File(dstDir))
	return nil
}
