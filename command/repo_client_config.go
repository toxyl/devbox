package command

import (
	"errors"
	"fmt"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/utils"
)

func RepoClientConfig(arg ...string) error {
	core.AppConfig.Repo.Client.Address = arg[0]
	core.AppConfig.Repo.Client.User = arg[1]
	core.AppConfig.Repo.Client.Password = utils.StringToSha256(arg[2])
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Repo client config updated!")
	return nil
}
