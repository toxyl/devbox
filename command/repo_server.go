package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
)

func RepoServer(arg ...string) error {
	srv := repo.NewServer(core.AppConfig.RepoAdminUser, core.AppConfig.RepoAdminPassword)
	return srv.ListenAndServe(core.AppConfig.RepoAddress, core.GetRepoDir())
}
