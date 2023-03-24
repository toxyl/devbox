package tmux

type TmuxWorkspace struct {
	created bool
	name    string
}

func (tw *TmuxWorkspace) add(cmd string) error {
	if !tw.created {
		if err := Exec("new-session", "-s", tw.name, "-d", cmd); err != nil {
			return err
		}
		tw.created = true
		return nil
	}
	return Exec("split-window", "-v", cmd)
}

func (tw *TmuxWorkspace) spawn(attach bool) error {
	if err := Exec("select-layout", "tiled"); err != nil {
		return err
	}
	if err := Exec("set-option", "-g", "mouse", "on"); err != nil {
		return err
	}
	if attach {
		if err := Exec("attach-session"); err != nil {
			return err
		}
	}
	return nil
}

func newTmuxWorkspace(name string) *TmuxWorkspace {
	tw := &TmuxWorkspace{
		name:    "devbox-" + name,
		created: false,
	}
	return tw
}

func SpawnWorkspace(name string, attach bool, commands ...string) error {
	if HasSession(name) {
		// we already have a session
		if !attach {
			// but we don't want to attach, let's quietly ignore this
			return nil
		}
		// let's attach
		return Exec("attach-session", "-t", "devbox-"+name)
	}

	tw := newTmuxWorkspace(name)
	for _, cmd := range commands {
		if err := tw.add(cmd); err != nil {
			return err
		}
	}
	return tw.spawn(attach)
}
