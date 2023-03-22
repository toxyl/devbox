package command

import (
	"fmt"
	"os"
)

func WorkspaceDestroy(arg ...string) error {
	name := arg[0]
	file := getWorkspaceConfigPath(name)
	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}
	path := getWorkspacePath(name)
	return os.RemoveAll(path)
}
