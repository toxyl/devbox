package sudo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Exec runs the given command with sudo privileges and returns its output.
func Exec(name string, arg ...string) ([]byte, error) {
	args := []string{"-S", "-p", "", "sh", "-c", fmt.Sprintf("%s %s", name, quoteArgs(arg))}
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = strings.NewReader("")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

// ExecInteractive runs the given command with sudo privileges and all streams attached.
func ExecInteractive(name string, arg ...string) error {
	args := []string{"-S", "-p", "", "sh", "-c", fmt.Sprintf("%s %s", name, quoteArgs(arg))}
	cmd := exec.Command("sudo", args...)
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

// quoteArgs quotes the given arguments so they can be passed as a single argument to a shell command.
func quoteArgs(args []string) string {
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			args[i] = fmt.Sprintf("'%s'", arg)
		}
	}
	return strings.Join(args, " ")
}
