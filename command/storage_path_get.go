package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func StoragePathGet(arg ...string) error {
	log.OK("Current storage directory: %s", glog.File(core.GetStorageDir()))
	return nil
}
