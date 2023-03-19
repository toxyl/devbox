package command

import (
	"strings"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func Enter(args ...string) error {
	devbox := args[0]
	log.Info("Entering %s", glog.Auto(devbox))
	_ = loadDevbox(devbox)
	log.Info(core.Config.Limits.String())
	run(
		devbox,
		SELF,
		SPAWN,
		devbox,
		core.Config.Env.Shell,
		"-c",
		strings.Join([]string{
			"echo \"" + core.Config.Commands.Start + "\"|" + core.Config.Env.Shell,
			core.Config.Env.Shell,
			"echo \"" + core.Config.Commands.Stop + "\"|" + core.Config.Env.Shell,
		}, " ; "),
	)
	return nil
}
