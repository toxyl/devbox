package command

import (
	"fmt"
	"os"
	"path/filepath"
)

func WorkspaceDestroy(arg ...string) error {
	path := arg[0]
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file := filepath.Join(path, ".workspace.yaml")
	if !fileExists(file) {
		return fmt.Errorf("workspace does not exist")
	}
	return os.RemoveAll(path)
}
