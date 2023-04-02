package command

import (
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/repo"
	"github.com/toxyl/glog"
)

func RepoDownload(fileNameRemote, fileNameLocal string) bool {
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		forceExit("Could not connect to repo server: "+glog.Error(err), core.EXIT_REPO_CONNECTION_FAILED)
	}
	isNew, err := client.DownloadFile(fileNameRemote, fileNameLocal, getStorageDir())
	if err != nil {
		forceExit("Could not download file from repo server: "+glog.Error(err), core.EXIT_REPO_DOWNLOAD_FAILED)
	}
	return isNew
}

func RepoUpload(fileNameSrc, fileNameDst string) {
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		forceExit("Could not connect to repo server: "+glog.Error(err), core.EXIT_REPO_CONNECTION_FAILED)
	}
	err = client.UploadFile(fileNameSrc, fileNameDst)
	if err != nil {
		forceExit("Could not upload file to repo server: "+glog.Error(err), core.EXIT_REPO_UPLOAD_FAILED)
	}
}
