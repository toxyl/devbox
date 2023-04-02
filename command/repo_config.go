package command

import (
	"errors"
	"fmt"

	"github.com/toxyl/devbox/core"
)

func RepoConfig(arg ...string) error {
	core.AppConfig.RepoAddress = arg[0]
	core.AppConfig.RepoPath = arg[1]
	core.AppConfig.RepoAdminUser = arg[2]
	core.AppConfig.RepoAdminPassword = arg[3]
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Credentials for repo administration updated!")
	return nil
}
