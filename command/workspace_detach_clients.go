package command

import (
	"github.com/toxyl/devbox/tmux"
)

func WorkspaceDetach(arg ...string) error {
	name := arg[0]
	return tmux.DetachClientsFromSession(name)
}
