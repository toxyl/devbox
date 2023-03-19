package tmux

import (
	"os"
	"os/exec"
)

type TmuxWorkspace struct {
	created bool
}

func (tw *TmuxWorkspace) exec(args ...string) error {
	cmd := exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Release()
		}
	}()
	return cmd.Run()
}

func (tw *TmuxWorkspace) Add(cmd string) error {
	if !tw.created {
		if err := tw.exec("new-session", "-d", cmd); err != nil {
			return err
		}
		tw.created = true
		return nil
	}
	return tw.exec("split-window", "-v", cmd)
}

func (tw *TmuxWorkspace) Spawn(attach bool) error {
	if err := tw.exec("select-layout", "tiled"); err != nil {
		return err
	}
	if err := tw.exec("set-option", "-g", "mouse", "on"); err != nil {
		return err
	}
	if attach {
		if err := tw.exec("attach-session"); err != nil {
			return err
		}
	}
	return nil
}

func NewTmuxWorkspace() *TmuxWorkspace {
	tw := &TmuxWorkspace{
		created: false,
	}
	return tw
}

func SpawnWorkspace(attach bool, commands ...string) error {
	tw := NewTmuxWorkspace()
	for _, cmd := range commands {
		if err := tw.Add(cmd); err != nil {
			return err
		}
	}
	return tw.Spawn(attach)
}
