package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func RepoServerInfo(arg ...string) error {
	log.Default(glog.Bold() + "SERVER" + glog.Reset())
	log.Default("Address:        %s", glog.Addr(core.AppConfig.Repo.Server.Address, false))
	log.Default("Path:           %s", glog.File(core.AppConfig.Repo.Server.Path))
	log.Default("Users:")
	for _, u := range core.AppConfig.Repo.Server.Users {
		log.Default("- User:         %s", glog.Highlight(u.Name))
		log.Default("  Password:     %s", glog.Highlight(u.Password))
		log.Default("  Admin:        %s", glog.Bool(u.Admin))
	}
	return nil
}
