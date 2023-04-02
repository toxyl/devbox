package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func RepoClientInfo(arg ...string) error {
	log.Default(glog.Bold() + "CLIENT" + glog.Reset())
	log.Default("Address:        %s", glog.Addr(core.AppConfig.Repo.Client.Address, false))
	log.Default("User:           %s", glog.Highlight(core.AppConfig.Repo.Client.User))
	log.Default("Password:       %s", glog.Highlight(core.AppConfig.Repo.Client.Password))
	return nil
}
