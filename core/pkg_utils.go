package core

import (
	"os"
	"path/filepath"
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

func GetStorageDir() string {
	return AppConfig.StoragePath
}

func GetRepoDir() string {
	return AppConfig.RepoPath
}

func GetWorkspaceDir() string {
	return filepath.Join(AppConfig.StoragePath, "workspace")
}

func GetWorkspaceConfigPath(name string) string {
	return filepath.Join(GetWorkspacePath(name), CONFIG_FILE)
}

func GetWorkspacePath(name string) string {
	return filepath.Join(GetWorkspaceDir(), name)
}

func GetDevboxDir() string {
	return filepath.Join(AppConfig.StoragePath, "devbox")
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
