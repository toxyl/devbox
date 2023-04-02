package command

import (
	"errors"
	"fmt"

	"github.com/toxyl/devbox/core"
)

func RepoCredentialsSet(arg ...string) error {
	core.AppConfig.AdminUser = arg[0]
	core.AppConfig.AdminPassword = arg[1]
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Credentials for repo administration updated!")
	return nil
}
