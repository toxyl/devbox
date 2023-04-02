package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func StoragePath(arg ...string) error {
	if len(arg) == 0 {
		log.OK("Current storage directory: %s", glog.File(core.GetStorageDir()))
		return nil
	}

	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !fileExists(path) {
		log.Question("The target directory %s does not exist, create it? [Y|n] "+glog.StoreCursor(), glog.File(path))
		time.Sleep(100 * time.Millisecond)
		fmt.Print(glog.RestoreCursor())
		var response string
		_, err := fmt.Scanln(&response)
		ok := true
		if err == nil {
			switch strings.ToLower(response) {
			case "y", "yes":
				ok = true
			case "n", "no":
				ok = false
			}
		}
		if !ok {
			return fmt.Errorf("user aborted creating target directory")
		}
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return errors.Join(err, fmt.Errorf("could not create target directory"))
		}
	}
	core.AppConfig.StoragePath = path
	if err := core.AppConfig.Save(); err != nil {
		return errors.Join(err, fmt.Errorf("could not save app config"))
	}

	UpdateBashCompletions()
	log.Success("Storage path updated!")
	log.Warning("Reload your shell using `exec bash` to refresh the Bash completions!")
	return nil
}
