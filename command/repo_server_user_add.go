package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/utils"
)

func RepoServerUserAdd(arg ...string) error {
	name := arg[0]
	password := arg[1]
	admin := false
	if strings.TrimSpace(strings.ToLower(arg[2])) == "true" {
		admin = true
	}
	add := true
	for i, u := range core.AppConfig.Repo.Server.Users {
		if u.Name == name {
			core.AppConfig.Repo.Server.Users[i].Password = utils.StringToSha256(password)
			core.AppConfig.Repo.Server.Users[i].Admin = admin
			add = false
			break
		}
	}
	if add {
		core.AppConfig.Repo.Server.Users = append(core.AppConfig.Repo.Server.Users, config.RepoUserConfig{
			Admin:    admin,
			Name:     name,
			Password: utils.StringToSha256(password),
		})
	}
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	log.Success("Repo server config updated!")
	log.Warning("You have to restart the repo server to apply the changes!")
	return nil
}
