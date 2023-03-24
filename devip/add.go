// Shameless rip of my own package at https://github.com/toxyl/devip
package devip

import (
	"net"
	"strings"

	"github.com/toxyl/devbox/sudo"
)

// add adds a loopback alias for the given IP.
func add(ip string) error {
	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		// we assume it didn't contain a CIDR, so it's a single machine
		ip += "/32"
	}
	output, err := sudo.Exec("ip", "address", "add", ip, "dev", "lo")
	if err != nil && !strings.Contains(string(output), "File exists") {
		return err
	}
	return nil
}

func Add(args ...string) {
	for _, alias := range args {
		add(alias)
	}
}
