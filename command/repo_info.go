package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func RepoInfo(arg ...string) error {
	glog.LoggerConfig.ShowIndicator = false
	glog.LoggerConfig.ShowRuntimeSeconds = false
	log.Default(glog.Bold() + "SERVER" + glog.Reset())
	log.Default("Address:        %s", glog.Addr(core.AppConfig.Repo.Server.Address, false))
	log.Default("Path:           %s", glog.File(core.AppConfig.Repo.Server.Path))
	log.Default("Users:")
	for _, u := range core.AppConfig.Repo.Server.Users {
		log.Default("- User:         %s", glog.Highlight(u.Name))
		log.Default("  Password:     %s", glog.Highlight(u.Password))
		log.Default("  Admin:        %s", glog.Bool(u.Admin))
	}
	log.Default("")
	log.Default(glog.Bold() + "CLIENT" + glog.Reset())
	log.Default("Address:        %s", glog.Addr(core.AppConfig.Repo.Client.Address, false))
	log.Default("User:           %s", glog.Highlight(core.AppConfig.Repo.Client.User))
	log.Default("Password:       %s", glog.Highlight(core.AppConfig.Repo.Client.Password))
	return nil
}
