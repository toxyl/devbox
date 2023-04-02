package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/devip"
	"github.com/toxyl/devbox/repo"
	"github.com/toxyl/devbox/tmux"
	"github.com/toxyl/glog"
)

func WorkspaceLaunch(arg ...string) error {
	name := arg[0]

	tarFileLocal := name + ".tar.gz"
	tarFileRemote := "workspace_" + name + ".tar.gz"

	log.Info("Checking %s...", glog.File(tarFileRemote))
	client := repo.NewClient(core.AppConfig.Repo.Client.User, core.AppConfig.Repo.Client.Password)
	err := client.Connect(core.AppConfig.Repo.Client.Address)
	if err != nil {
		forceExit("Could not connect to repo server: "+glog.Error(err), core.EXIT_REPO_CONNECTION_FAILED)
	}
	uptodate, err := client.CheckIfFileIsUpToDate(tarFileLocal, tarFileRemote, getStorageDir())
	if err != nil {
		if err.Error() == "ERROR File not found" {
			uptodate = true
		} else {
			forceExit("Could not check if file is up-to-date: "+glog.Error(err), core.EXIT_REPO_CHECK_FAILED)
		}
	}
	if !uptodate {
		log.Question("There is a newer workspace version, update? [y|N] " + glog.StoreCursor())
		time.Sleep(100 * time.Millisecond)
		fmt.Print(glog.RestoreCursor())
		var response string
		_, err := fmt.Scanln(&response)
		ok := false
		if err == nil {
			switch strings.ToLower(response) {
			case "y", "yes":
				ok = true
			case "n", "no":
				ok = false
			default:
				ok = false
			}
		}
		if ok {
			err = WorkspacePull(name)
			if err != nil {
				forceExit("Could not download file from repo server: "+glog.Error(err), core.EXIT_REPO_DOWNLOAD_FAILED)
			}
		}
	}

	file := getWorkspaceConfigPath(name)

	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}

	w, err := config.OpenWorkspace(file)
	if err != nil {
		return err
	}
	commands := make([]string, len(w.Devboxes))
	for i, devbox := range w.Devboxes {
		dboxPath := getDevboxConfigPath(devbox.Name)
		if !fileExists(dboxPath) {
			forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(dboxPath)), core.EXIT_DEVBOX_NOT_FOUND)
		}
		must(config.DevboxToFile(&devbox.Config, dboxPath))
		arg := ""
		if devbox.Delay > 0 {
			arg = "sleep " + fmt.Sprint(devbox.Delay) + " ; "
		}
		commands[i] = arg + core.APP_NAME + " " + ENTER + " " + devbox.Name
	}

	for _, ip := range w.IPs {
		devip.Add(ip)
	}

	defer func() {
		for _, ip := range w.IPs {
			devip.Remove(ip)
		}
	}()

	err = tmux.SpawnWorkspace(name, true, commands...)
	// store current container config back into workspace config
	for i, devbox := range w.Devboxes {
		dboxPath := getDevboxConfigPath(devbox.Name)
		if !fileExists(dboxPath) {
			forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(dboxPath)), core.EXIT_DEVBOX_NOT_FOUND)
		}
		if err := devbox.Config.Load(dboxPath); err != nil {
			return err
		}
		w.Devboxes[i] = devbox
	}

	w.Save(file)

	return err
}
