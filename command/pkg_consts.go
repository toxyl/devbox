package command

const (
	SPAWN               = "spawn"
	MAKE                = "make"
	ENTER               = "enter"
	EXEC                = "exec"
	STORE               = "store"
	DESTROY             = "destroy"
	WORKSPACE_LAUNCH    = "workspace-launch"
	WORKSPACE_ADD       = "workspace-add"
	WORKSPACE_REMOVE    = "workspace-remove"
	WORKSPACE_IP_ADD    = "workspace-ip-add"
	WORKSPACE_IP_REMOVE = "workspace-ip-remove"
	WORKSPACE_STORE     = "workspace-store"
	WORKSPACE_RESTORE   = "workspace-restore"
	WORKSPACE_DESTROY   = "workspace-destroy"
	WORKSPACE_DETACH    = "workspace-detach"
	STORAGE_PATH_GET    = "storage-path-get"
	STORAGE_PATH_SET    = "storage-path-set"
	SELF                = "/proc/self/exe"
)

const (
	TEXT_SUDO_REQUIRED = "This application requires superuser privileges to run.\nPlease restart the application with 'sudo'."
)
