package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
	"github.com/toxyl/glog"
)

func RepoServer(arg ...string) error {
	glog.LoggerConfig.ShowRuntimeSeconds = true
	srv := repo.NewServer(core.AppConfig.Repo.Server.Path)
	for _, u := range core.AppConfig.Repo.Server.Users {
		if u.Admin {
			srv.AddAdmin(u.Name, u.Password)
			continue
		}
		srv.AddUser(u.Name, u.Password)
	}
	return srv.ListenAndServe(core.AppConfig.Repo.Server.Address)
}
