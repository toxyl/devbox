package command

import (
	"os"
	"os/exec"

	"github.com/toxyl/devbox/core"
)

func Spawn(args ...string) error {
	devbox := args[0]
	path := loadDevbox(devbox)
	args = args[1:]

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = core.Config.GetEnv()
	bindAll(devbox)
	setHostname(devbox)
	pivotRoot(path)
	defer func() {
		unbindAll()
		unmountSpecialFS()
		if cmd.Process != nil {
			cmd.Process.Release()
		}
	}()

	must(cmd.Run(), "device or resource busy")
	return nil
}
