package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/toxyl/devbox/cgroups"
	"github.com/toxyl/devbox/core"
)

func run(devbox string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	_, memHardMax, memSoftMax, memSwap := core.Config.Limits.GetMemoryLimits()
	cgpath, err := cgroups.CreateCGroup(
		core.APP_NAME, devbox,
		core.Config.Limits.GetCPUQuota(),
		memHardMax, memSoftMax, memSwap,
		core.Config.Limits.PIDs)
	must(err)
	defer func() {
		cgroups.DestroyCGroup(core.APP_NAME, devbox)
		if cmd.Process != nil {
			cmd.Process.Release()
		}
	}()
	fd, err := os.Open(cgpath)
	must(err)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		GidMappingsEnableSetgroups: true,
		UseCgroupFD:                true,
		CgroupFD:                   int(fd.Fd()),
	}
	if core.Config.Options.MapUsersAndGroups {
		cmd.SysProcAttr.UidMappings = getUIDMappings()
		cmd.SysProcAttr.GidMappings = getGIDMappings()
	} else {
		cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			}, {
				ContainerID: 1,
				HostID:      100000,
				Size:        65534,
			},
		}
		cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			}, {
				ContainerID: 1,
				HostID:      100000,
				Size:        65534,
			},
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	must(cmd.Run(), "device or resource busy")
}

func Run() {
	cmd := findCommand(os.Args[1])
	if cmd != nil {
		checkError(cmd.Run())
	}
}
