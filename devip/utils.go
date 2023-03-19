// Shameless rip of my own package at https://github.com/toxyl/devip
package devip

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

func getList() map[string]string {
	res := map[string]string{}
	output, err := sudoExec("ip", "address", "show", "dev", "lo")
	if err != nil {
		return res
	}
	re := regexp.MustCompile(`inet\s+(\d+\.\d+\.\d+\.\d+)/(\d+)`)
	matches := re.FindAllStringSubmatch(string(output), -1)
	ips := make([]net.IP, 0, len(matches))
	cidrs := make([]string, 0, len(matches))
	for _, match := range matches {
		ip := net.ParseIP(match[1])
		if ip != nil && !isLocalhost(ip.String()) {
			ips = append(ips, ip)
			cidrs = append(cidrs, match[2])
		}
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	for i, ip := range ips {
		res[ip.String()] = fmt.Sprintf("%s/%s", ip.String(), cidrs[i])
	}
	return res
}

// sudoExec runs the given command with sudo privileges and returns its output.
func sudoExec(name string, arg ...string) ([]byte, error) {
	args := []string{"-S", "-p", "", "sh", "-c", fmt.Sprintf("%s %s", name, quoteArgs(arg))}
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = strings.NewReader("")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
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

// isLocalhost checks whether the given alias is a loopback address.
func isLocalhost(alias string) bool {
	ip := net.ParseIP(alias)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
