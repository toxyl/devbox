package main

import (
	"fmt"
	"os"

	"github.com/toxyl/devbox/command"
	"github.com/toxyl/devbox/core"
	"github.com/toxyl/glog"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Printf("%s\n", command.TEXT_SUDO_REQUIRED)
		os.Exit(core.EXIT_NEED_SUDO)
	}

	glog.LoggerConfig.ShowSubsystem = false
	glog.LoggerConfig.ShowDateTime = false
	glog.LoggerConfig.ShowRuntimeMilliseconds = false
	glog.LoggerConfig.ShowRuntimeSeconds = true

	core.InitErrorRegistry()

	core.RegisterCommand(
		command.MAKE,
		"Creates a new `devbox` from a `tarball` (either filepath or URL).\nIf `tarball` is a local file it will be unpacked to the devbox location.\nIf `tarball`` is a URL the file will be downloaded first and then unpacked.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DEVBOX_LIST,
				Name:    "devbox",
				Example: "kinetic",
			},
			{
				Type:    core.ARG_TYPE_TARBALL,
				Name:    "tarball",
				Example: "https://cloud-images.ubuntu.com/minimal/releases/kinetic/release/ubuntu-22.10-minimal-cloudimg-amd64-root.tar.xz",
			},
		},
		command.Make,
	)

	core.RegisterCommand(
		command.ENTER,
		"Enters the `devbox`.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DEVBOX_LIST,
				Name:    "devbox",
				Example: "kinetic",
			},
		},
		command.Enter,
	)

	core.RegisterCommand(
		command.EXEC,
		"Executes `command` in the `devbox`.\n"+glog.HighlightWarning("Wrap multiple arguments in quotes!"),
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DEVBOX_LIST,
				Name:    "devbox",
				Example: "kinetic",
			},
			{
				Type:    core.ARG_TYPE_COMMAND,
				Name:    "command",
				Example: "'hostname ; hostname justanotherct ; hostname'",
			},
		},
		command.Exec,
	)

	core.RegisterCommand(
		command.STORE,
		"Stores the `devbox` as `tarball`.\nSupported extensions: .tar.gz, .tar.xz",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DEVBOX_LIST,
				Name:    "devbox",
				Example: "kinetic",
			},
			{
				Type:    core.ARG_TYPE_TARBALL,
				Name:    "tarball",
				Example: "kinetic.tar.gz",
			},
		},
		command.Store,
	)

	core.RegisterCommand(
		command.DESTROY,
		"Removes the `devbox`.\nIt "+glog.Bold()+glog.WrapRed("doesn't")+glog.Reset()+" ask for confirmation!",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DEVBOX_LIST,
				Name:    "devbox",
				Example: "kinetic",
			},
		},
		command.Destroy,
	)

	core.RegisterCommand(
		command.WORKSPACE_ADD,
		"Adds `devbox`es to the workspace in `path`.\nIf the workspace doesn't exist a new one will be created.\nEach argument can have an optional `:delay` parameter which is used to set the startup delay for the `devbox` in seconds.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Variadic:       true,
				Type:           core.ARG_TYPE_DEVBOX_LIST,
				Name:           "devbox",
				Example:        "bionic kinetic:10 focal",
				OptionalSuffix: ":delay",
			},
		},
		command.WorkspaceAdd,
	)

	core.RegisterCommand(
		command.WORKSPACE_REMOVE,
		"Removes `devbox`es from the workspace in `path`.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Variadic:       true,
				Type:           core.ARG_TYPE_DEVBOX_LIST,
				Name:           "devbox",
				Example:        "bionic kinetic focal",
				OptionalSuffix: "",
			},
		},
		command.WorkspaceRemove,
	)

	core.RegisterCommand(
		command.WORKSPACE_IP_ADD,
		"Adds `IP`s to the workspace in `path`.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Variadic:       true,
				Type:           core.ARG_TYPE_IP,
				Name:           "IP",
				Example:        "192.168.0.1 192.168.0.2 1.0.0.0/24",
				OptionalSuffix: "/prefix",
			},
		},
		command.WorkspaceIPAdd,
	)

	core.RegisterCommand(
		command.WORKSPACE_IP_REMOVE,
		"Removes `IP`s from the workspace in `path`.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Variadic:       true,
				Type:           core.ARG_TYPE_IP,
				Name:           "IP",
				Example:        "192.168.0.1 192.168.0.2 1.0.0.0/24",
				OptionalSuffix: "/prefix",
			},
		},
		command.WorkspaceIPRemove,
	)

	core.RegisterCommand(
		command.WORKSPACE_LAUNCH,
		"Launches the workspace in `path` in a tmux session.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
		},
		command.WorkspaceLaunch,
	)

	core.RegisterCommand(
		command.WORKSPACE_STORE,
		"Stores the workspace in `path` as tarball.\nSupported extensions: .tar.gz, .tar.xz",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Type:    core.ARG_TYPE_TARBALL,
				Name:    "tarball",
				Example: "~/my-workspace.tar.gz",
			},
		},
		command.WorkspaceStore,
	)

	core.RegisterCommand(
		command.WORKSPACE_RESTORE,
		"Restores the workspace from `tarball` to `path`.\n"+glog.Bold()+glog.WrapOrange("Warning:")+glog.Reset()+" Existing devboxes will be overwritten!",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
			{
				Type:    core.ARG_TYPE_TARBALL,
				Name:    "tarball",
				Example: "~/my-workspace.tar.gz",
			},
		},
		command.WorkspaceRestore,
	)

	core.RegisterCommand(
		command.WORKSPACE_DESTROY,
		"Completely removes the workspace in `path`.",
		core.ArgInfoList{
			{
				Type:    core.ARG_TYPE_DIR,
				Name:    "path",
				Example: "/tmp/my-workspace/",
			},
		},
		command.WorkspaceDestroy,
	)

	// Hidden Commands
	// These are not shown in help texts and used internally.
	core.RegisterHiddenCommand(
		command.SPAWN,
		core.ArgInfoList{
			{
				Variadic: true,
				Type:     core.ARG_TYPE_FILE,
				Name:     "command",
			},
		},
		command.Spawn,
	)

	if len(os.Args) < 2 {
		fmt.Println("\n" +
			glog.Bold() + glog.WrapYellow("Welcome to Devbox!") + glog.Reset() +
			"\n" +
			"\n" +
			glog.Underline() + "Available Commands" + glog.Reset() +
			"\n" +
			"(omit args to see help)" +
			"\n")
		core.ListCommands()
		command.UpdateBashCompletions()
		fmt.Println()
		os.Exit(core.EXIT_OK)
	}

	command.Run()
}