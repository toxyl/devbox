package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/toxyl/devbox/bashcompletion"
	"github.com/toxyl/devbox/config"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func loadDevbox(name string) string {
	p := getDevboxPath(name)
	if !fileExists(p) {
		forceExit(fmt.Sprintf("Devbox %s does not exist!", glog.Auto(name)), core.EXIT_DEVBOX_NOT_FOUND)
	}
	createMissingDevboxFiles(name)
	c, err := config.DevboxFromFile(getDevboxConfigPath(name))
	must(err)
	core.Config = &c
	return getDevboxPath(name)
}

func bind(devbox, src, dst string) {
	target := filepath.Join(getDevboxPath(devbox), strings.TrimLeft(src, "/"))
	source, _ := filepath.Abs(dst)
	log.Info("Binding %s to %s:%s", glog.File(source), glog.Auto(devbox), glog.File(src))
	_ = os.RemoveAll(target) // if target exist, remove it!
	must(os.MkdirAll(target, 0755))
	must(os.MkdirAll(source, 0755))
	must(syscall.Mount(source, target, "bind", syscall.MS_BIND, ""))
}

func undbind(src string) {
	must(syscall.Unmount("/"+strings.TrimLeft(src, "/"), 0), "device or resource busy")
}

func mountSpecialFS() {
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""))
	must(syscall.Mount("devpts", "dev/pts", "devpts", 0, ""))

}

func unmountSpecialFS() {
	must(syscall.Unmount("/proc", 0), "device or resource busy")
	must(syscall.Unmount("/tmp", 0), "device or resource busy")
	must(syscall.Unmount("/dev/pts", 0), "device or resource busy")
}

func pivotRoot(path string) {
	must(syscall.Chroot(path))
	must(syscall.Chdir("/"))
	mountSpecialFS()

}

func bindAll(devbox string) {
	for bindSrc, bindDst := range core.Config.Binds {
		if bindSrc == "" || bindDst == "" {
			continue
		}
		bind(devbox, bindSrc, bindDst)
	}
	path := getDevboxPath(devbox)
	if core.Config.Options.BindAll {
		bind(devbox, "/devboxes", filepath.Dir(path))
	} else {
		_ = os.Remove(filepath.Join(path, "devboxes")) // make sure the dir is gone if the setting has changed
	}
}

func unbindAll() {
	for bindSrc := range core.Config.Binds {
		undbind(bindSrc)
	}
}

func setHostname(name string) {
	must(syscall.Sethostname([]byte(name)))
}

func getUsersMap() map[string]int {
	users := make(map[string]int)

	// Get the current user's UID.
	currentUID := os.Getuid()

	// Open the password file.
	file, err := os.Open("/etc/passwd")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read each line in the file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Parse the user's information.
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 3 {
			continue
		}

		// Parse the UID and user name.
		uid, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		username := fields[0]

		// Add the user to the map if their UID is not the current user's UID.
		if uid != currentUID {
			users[username] = uid
		}
	}

	// Check for any errors while scanning the file.
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return users
}

func getUIDMappings() []syscall.SysProcIDMap {
	mappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		},
	}
	highestUID := 0
	for _, uid := range getUsersMap() {
		if uid == 0 {
			continue // skip root user
		}
		mapping := syscall.SysProcIDMap{
			ContainerID: uid,
			HostID:      uid,
			Size:        1,
		}
		mappings = append(mappings, mapping)
		if uid > highestUID {
			highestUID = uid
		}
	}
	mappings = append(mappings, syscall.SysProcIDMap{
		ContainerID: highestUID + 1,
		HostID:      100000,
		Size:        65534,
	})
	return mappings
}

func getGroupsMap() map[string]int {
	groups := make(map[string]int)

	// Get the current user's GID.
	currentGID := os.Getgid()

	// Open the group file.
	file, err := os.Open("/etc/group")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read each line in the file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Parse the group's information.
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 3 {
			continue
		}

		// Parse the GID and group name.
		gid, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		groupname := fields[0]

		// Add the group to the map if their GID is not the current user's GID.
		if gid != currentGID {
			groups[groupname] = gid
		}
	}

	// Check for any errors while scanning the file.
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return groups
}

func getGIDMappings() []syscall.SysProcIDMap {
	mappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		},
	}
	highestGID := 0
	for _, gid := range getGroupsMap() {
		if gid == 0 {
			continue // skip root group
		}
		mapping := syscall.SysProcIDMap{
			ContainerID: gid,
			HostID:      gid,
			Size:        1,
		}
		mappings = append(mappings, mapping)
		if gid > highestGID {
			highestGID = gid
		}
	}
	mappings = append(mappings, syscall.SysProcIDMap{
		ContainerID: highestGID + 1,
		HostID:      100000,
		Size:        65534,
	})
	return mappings
}

func UpdateBashCompletions() {
	core.Must(bashcompletion.Make(
		core.APP_NAME,
		core.GetDevboxDir(),
		core.GetCommandNames(),
		core.GetCommandData(),
	))
}
