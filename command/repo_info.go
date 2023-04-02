package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func RepoInfo(arg ...string) error {
	log.Default("Address:        %s", glog.Addr(core.AppConfig.RepoAddress, false))
	log.Default("Path:           %s", glog.File(core.GetRepoDir()))
	log.Default("Admin User:     %s", glog.Highlight(core.AppConfig.RepoAdminUser))
	log.Default("Admin Password: %s", glog.Highlight(utils.StringToSha256(core.AppConfig.RepoAdminPassword)))
	return nil
}
