package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
)

func RepoServer(arg ...string) error {
	srv := repo.NewServer(core.AppConfig.AdminUser, core.AppConfig.AdminPassword)
	return srv.ListenAndServe(arg[0], arg[1])
}
