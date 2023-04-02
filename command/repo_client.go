package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
	"github.com/toxyl/glog"
)

func RepoDownload(arg ...string) error {
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		log.Error("Could not connect to repo server: %s", glog.Error(err))
		return nil
	}
	err = client.DownloadFile(arg[0])
	if err != nil {
		log.Error("Could not download file from repo server: %s", glog.Error(err))
	}
	return nil
}

func RepoUpload(arg ...string) error {
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		log.Error("Could not connect to repo server: %s", glog.Error(err))
		return nil
	}
	err = client.UploadFile(arg[0], arg[1])
	if err != nil {
		log.Error("Could not upload file to repo server: %s", glog.Error(err))
	}
	return nil
}
