package command

import (
	"github.com/toxyl/devbox/core"
)

func Exec(args ...string) error {
	devbox := args[0]
	command := args[1]
	_ = loadDevbox(devbox)
	run(
		devbox,
		SELF,
		SPAWN,
		devbox,
		core.Config.Env.Shell,
		"-c",
		command,
	)
	return nil
}
