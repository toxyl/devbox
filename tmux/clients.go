package tmux

import (
	"strings"

	"github.com/toxyl/devbox/sudo"
)

func GetClients() (map[string][]string, error) {
	clients := map[string][]string{}
	output, err := sudo.Exec("tmux", "list-clients")
	if err != nil && !strings.Contains(string(output), "no server running on") {
		return clients, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) == 0 {
		return clients, nil
	}

	for _, ln := range lines {
		lnStr := string(ln)
		if lnStr == "" {
			continue
		}
		elems := strings.Split(lnStr, " ")
		clientID := elems[0]
		sessionID := elems[1]
		if strings.HasPrefix(sessionID, "devbox-") {
			if _, ok := clients[sessionID]; !ok {
				clients[sessionID] = []string{}
			}
			clients[sessionID] = append(clients[sessionID], clientID[0:len(clientID)-1])
		}
	}

	return clients, nil
}
