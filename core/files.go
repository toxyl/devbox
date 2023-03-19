package core

import (
	"os"
)

func createConfigIfMissing(name string) {
	if !FileExists(GetDevboxConfigPath(name)) {
		Must(os.WriteFile(GetDevboxConfigPath(name), []byte(defaultConfig), 0644))
	}
}

func createStartScriptIfMissing(name string) {
	if !FileExists(getDevboxStartScriptPath(name)) {
		Must(os.WriteFile(getDevboxStartScriptPath(name), []byte(defaultStartScript), 0755))
	}
}

func createStopScriptIfMissing(name string) {
	if !FileExists(getDevboxStopScriptPath(name)) {
		Must(os.WriteFile(getDevboxStopScriptPath(name), []byte(defaultStopScript), 0755))
	}
}

func CreateMissingDevboxFiles(name string) {
	createConfigIfMissing(name)
	createStartScriptIfMissing(name)
	createStopScriptIfMissing(name)
}
