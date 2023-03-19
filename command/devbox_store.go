package command

import (
	"fmt"
	"path/filepath"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func Store(arg ...string) error {
	name := arg[0]
	src := getDevboxPath(name)
	dst, err := filepath.Abs(arg[1])
	must(err)
	if !fileExists(src) {
		forceExit(fmt.Sprintf("Devbox %s does not exist in %s!", glog.Auto(name), glog.File(src)), core.EXIT_FILE_NOT_FOUND)
	}
	log.Info("Storing devbox %s...", glog.Auto(name))
	must(tar.FromDir(src, dst))
	log.Success("Stored to %s", glog.File(dst))
	return err
}
