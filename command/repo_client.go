package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
)

func RepoDownload(arg ...string) error {
	client := repo.NewClient(core.AppConfig.AdminUser, core.AppConfig.AdminPassword)
	err := client.Connect(arg[0])
	if err != nil {
		return err
	}
	return client.DownloadFile(arg[1])
}

func RepoUpload(arg ...string) error {
	client := repo.NewClient(core.AppConfig.AdminUser, core.AppConfig.AdminPassword)
	err := client.Connect(arg[0])
	if err != nil {
		return err
	}
	return client.UploadFile(arg[1], arg[2])
}
