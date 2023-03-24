package core

import (
	_ "embed"
	"os"
	"regexp"

	"github.com/toxyl/devbox/config"
	"github.com/toxyl/glog"
)

//go:embed resources/config.yaml
var defaultConfig string

//go:embed resources/start
var defaultStartScript string

//go:embed resources/stop
var defaultStopScript string

var (
	reToken     = regexp.MustCompile(`\{[^\}]+\}|\[[^\]]+\]|<[^>]+>`)
	log         = glog.NewLoggerSimple("core")
	Config      = &config.Config{}
	cmdReg      = []*command{}
	ERRORS      = map[string]error{}
	FATALS      = map[string]error{}
	fatalErrors = map[string]int{
		ERR_INVALID_ARGS:   EXIT_INVALID_ARGS,
		ERR_MISSING_ARGS:   EXIT_MISSING_ARGS,
		ERR_TOO_MANY_ARGS:  EXIT_TOO_MANY_ARGS,
		ERR_FILE_NOT_FOUND: EXIT_FILE_NOT_FOUND,
		ERR_FATAL:          EXIT_UNKNOWN_FATAL,
	}
	errReg      = glog.NewGErrorRegistry()
	storagePath = os.TempDir()
)
