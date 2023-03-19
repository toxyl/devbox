// Shameless rip of my own package at https://github.com/toxyl/devip
package devip

import (
	"fmt"
	"net"
	"strings"
)

// remove removes the loopback alias of the given IP.
func remove(ip string) error {
	if isLocalhost(ip) {
		return nil
	}

	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		// we assume it didn't contain a CIDR, so it's a single machine
		ip += "/32"
	}

	output, err := sudoExec("ip", "address", "del", ip, "dev", "lo")
	if err != nil {
		if !strings.Contains(string(output), "Cannot assign requested address") {
			return fmt.Errorf(string(output) + ": " + err.Error())
		}
	}
	return nil
}

func Remove(args ...string) {
	if len(args) == 1 && args[0] == "all" {
		for _, cidr := range getList() {
			remove(cidr)
		}
		return
	}
	for _, alias := range args {
		remove(alias)
	}
}
