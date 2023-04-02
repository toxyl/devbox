package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
	"github.com/toxyl/glog"
)

func RepoDownload(arg ...string) error {
	log.Info("Downloading %s to %s...", glog.File(arg[0]), glog.File(arg[1]))
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		forceExit("Could not connect to repo server: "+glog.Error(err), core.EXIT_REPO_CONNECTION_FAILED)
	}
	err = client.DownloadFile(arg[0], arg[1], getStorageDir())
	if err != nil {
		forceExit("Could not download file from repo server: "+glog.Error(err), core.EXIT_REPO_DOWNLOAD_FAILED)
	}
	return nil
}

func RepoUpload(arg ...string) error {
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		forceExit("Could not connect to repo server: "+glog.Error(err), core.EXIT_REPO_CONNECTION_FAILED)
	}
	err = client.UploadFile(arg[0], arg[1])
	if err != nil {
		forceExit("Could not upload file to repo server: "+glog.Error(err), core.EXIT_WORKSPACE_PUSH_FAILED)
	}
	return nil
}
