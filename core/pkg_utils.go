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

func GetDevboxDir() string {
	return filepath.Join(os.TempDir(), APP_NAME)
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
