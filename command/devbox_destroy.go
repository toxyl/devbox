package command

import (
	"fmt"
	"os"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func Destroy(arg ...string) error {
	name := arg[0]
	path := getDevboxPath(name)
	if !fileExists(path) {
		forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(name)), core.EXIT_DEVBOX_NOT_FOUND)
	}
	log.Info("Destroying devbox %s (%s)...", glog.Auto(name), glog.File(path))
	must(os.RemoveAll(path))
	log.Success("And, %s is gone!", glog.Auto(name))
	return nil
}
