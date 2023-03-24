package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/glog"
)

func FileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}

func SetStorageDir(path string) error {
	if !FileExists(path) {
		return fmt.Errorf("target directory does not exist: %s", glog.File(path))
	}
	storagePath = path
	return nil
}

func GetStorageDir() string {
	return storagePath
}

func GetWorkspaceDir() string {
	return filepath.Join(storagePath, "workspace")
}

func GetWorkspaceConfigPath(name string) string {
	return filepath.Join(GetWorkspacePath(name), CONFIG_FILE)
}

func GetWorkspacePath(name string) string {
	return filepath.Join(GetWorkspaceDir(), name)
}

func GetDevboxDir() string {
	return filepath.Join(storagePath, "devbox")
}

func GetDevboxPath(name string) string {
	return filepath.Join(GetDevboxDir(), name)
}

func GetDevboxConfigPath(name string) string {
	return filepath.Join(GetDevboxPath(name), CONFIG_FILE)
}

func getDevboxStartScriptPath(name string) string {
	return filepath.Join(GetDevboxPath(name), START_SCRIPT_FILE)
}

func getDevboxStopScriptPath(name string) string {
	return filepath.Join(GetDevboxPath(name), STOP_SCRIPT_FILE)
}
