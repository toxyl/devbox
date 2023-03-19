package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/toxyl/glog"
)

func InitErrorRegistry() {
	ERRORS = map[string]error{}
	FATALS = map[string]error{}
	for msg, exit := range fatalErrors {
		FATALS[msg] = errors.New(msg)
		errReg.Register(FATALS[msg], true, exit)
	}
}

// CheckError will try to register the given error as a non-fatal error
// and then check if the given error matches a known one (which should always be the case).
// If err == nil the function will exit immediately.
func CheckError(err error) {
	if err == nil {
		return
	}
	if dbusErr, ok := err.(dbus.Error); ok {
		err = fmt.Errorf(dbusErr.Error())
	}
	errReg.Register(err, false, EXIT_OK) // make sure we treat it as an error from now on, but don';'t fail
	errReg.Check(err, err.Error())
}

// CheckFatal will check if the given error is a fatal error and if so, it will exit with the registered exit code.
// If err == nil the function will exit immediately.
func CheckFatal(err error) {
	if err == nil {
		return
	}
	errReg.Check(err, err.Error())
}

func ForceFatal(msg string) {
	FATALS[msg] = errors.New(msg)
	errReg.Register(FATALS[msg], true, EXIT_UNKNOWN_FATAL)
	errReg.Check(FATALS[msg], glog.Bold()+"FATAL"+glog.Reset())
}

func Must(err error, exclude ...string) {
	if err == nil {
		return
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	for _, e := range exclude {
		if strings.Contains(msg, strings.ToLower(strings.TrimSpace(e))) {
			return
		}
	}
	ForceFatal(err.Error())
}

func ForceExit(msg string, exitCode int) {
	log.Error(msg)
	os.Exit(exitCode)
}
