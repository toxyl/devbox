package command

import (
	"errors"
	"fmt"

	"github.com/toxyl/devbox/core"
)

func RepoServerConfig(arg ...string) error {
	core.AppConfig.Repo.Server.Address = arg[0]
	core.AppConfig.Repo.Server.Path = arg[1]
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Repo server config updated!")
	return nil
}
