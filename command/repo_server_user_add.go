package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
)

func RepoServerUserAdd(arg ...string) error {
	name := arg[0]
	password := arg[1]
	admin := false
	if len(arg) == 3 && strings.TrimSpace(strings.ToLower(arg[2])) == "true" {
		admin = true
	}
	core.AppConfig.Repo.Server.Users = append(core.AppConfig.Repo.Server.Users, config.RepoUserConfig{
		Admin:    admin,
		Name:     name,
		Password: password,
	})
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Repo server config updated!")
	log.Warning("You have to restart the repo server to apply the changes!")
	return nil
}
