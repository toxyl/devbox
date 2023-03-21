package command

import (
	"fmt"
	"path/filepath"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func WorkspaceRestore(arg ...string) error {
	tarFile := arg[1]
	tarFile, err := filepath.Abs(tarFile)
	if err != nil {
		return err
	}
	if !fileExists(tarFile) {
		return fmt.Errorf("the source file %s does not exist", tarFile)
	}

	dstDir := arg[0]
	dstDir, err = filepath.Abs(dstDir)
	if err != nil {
		return err
	}
	log.Success("Restoring workspace from %s", glog.File(tarFile))

	err = tar.ToDir(tarFile, dstDir)
	if err != nil {
		return err
	}

	file := filepath.Join(dstDir, ".workspace.yaml")
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
	log.Success("Restored workspace to %s", glog.File(dstDir))
	return nil
}
