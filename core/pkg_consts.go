package core

const (
	EXIT_OK = iota
	EXIT_NEED_SUDO
	EXIT_INVALID_ARGS
	EXIT_MISSING_ARGS
	EXIT_TOO_MANY_ARGS
	EXIT_FILE_NOT_FOUND
	EXIT_DEVBOX_NOT_FOUND
	EXIT_UNKNOWN_ERROR
	EXIT_UNKNOWN_FATAL
)

const (
	ERR_INVALID_ARGS     = "invalid args"
	ERR_MISSING_ARGS     = "missing required argument"
	ERR_TOO_MANY_ARGS    = "too many arguments"
	ERR_FILE_NOT_FOUND   = "no such file or directory"
	ERR_DEVBOX_NOT_FOUND = "devbox not found"
	ERR_FATAL            = "fatal error"
	CONFIG_FILE          = "config.yaml"
	START_SCRIPT_FILE    = "usr/local/bin/start"
	STOP_SCRIPT_FILE     = "usr/local/bin/stop"
	APP_NAME             = "devbox"
)
