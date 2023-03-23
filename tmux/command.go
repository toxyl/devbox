package tmux

import (
	"github.com/toxyl/devbox/sudo"
)

func Exec(args ...string) error {
	return sudo.ExecInteractive("tmux", args...)
}
