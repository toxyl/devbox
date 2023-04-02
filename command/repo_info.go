package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func RepoInfo(arg ...string) error {
	log.Default("Address:        %s", glog.Addr(core.AppConfig.Repo.Server.Address, false))
	log.Default("Path:           %s", glog.File(core.AppConfig.Repo.Server.Path))
	log.Default("Users:")
	for _, u := range core.AppConfig.Repo.Server.Users {
		log.Default("- User:         %s", glog.Highlight(u.Name))
		log.Default("  Admin:        %s", glog.Bool(u.Admin))
		log.Default("  Password:     %s", glog.Highlight(utils.StringToSha256(u.Password)))
	}
	return nil
}
