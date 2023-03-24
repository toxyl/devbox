package tmux

import (
	"strings"

	"github.com/toxyl/devbox/sudo"
)

func GetSessions() ([]string, error) {
	sessions := []string{}
	output, err := sudo.Exec("tmux", "list-sessions")
	if err != nil {
		return sessions, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) == 0 {
		return sessions, nil
	}

	for _, ln := range lines {
		lnStr := string(ln)
		if lnStr == "" {
			continue
		}
		id := strings.Split(lnStr, ":")[0]
		if strings.HasPrefix(id, "devbox-") {
			sessions = append(sessions, id)
		}
	}

	return sessions, nil
}

func HasSession(name string) bool {
	sess, err := GetSessions()
	if err != nil {
		return false
	}
	name = "devbox-" + name
	for _, session := range sess {
		if session == name {
			return true
		}
	}
	return false
}

func DetachClientsFromSession(name string) error {
	clients, err := GetClients()
	if err != nil {
		return err
	}
	name = "devbox-" + name
	if _, ok := clients[name]; !ok {
		// probably session is not active anymore, let's ignore it
		return nil
	}
	for _, clientID := range clients[name] {
		if err := Exec("detach-client", "-t", clientID); err != nil {
			return err
		}
	}
	return nil
}
