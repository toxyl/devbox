package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toxyl/devbox/core"
	"github.com/toxyl/devbox/downloader"
	"github.com/toxyl/devbox/tar"
	"github.com/toxyl/glog"
)

func unpack(arg ...string) error {
	name := arg[0]
	dst := getDevboxPath(name)
	src, err := filepath.Abs(arg[1])
	must(err)
	if !fileExists(src) {
		forceExit(fmt.Sprintf("Source %s does not exist!", glog.File(src)), core.EXIT_FILE_NOT_FOUND)
	}
	must(tar.ToDir(src, dst))
	createMissingDevboxFiles(name)
	return err
}

func Make(arg ...string) error {
	name := arg[0]
	dst := getDevboxPath(name)
	src, err := filepath.Abs(arg[1])
	must(err)

	if glog.IsURL(arg[1]) {
		_ = os.MkdirAll(dst, 0777)
		fdst := dst + "/../"
		f := arg[1]
		log.Info("Downloading %s to %s", glog.Auto(f), glog.Auto(fdst))
		f, err := downloader.Download(f, fdst)
		must(err)
		src = f
		arg[1] = f
	} else if !glog.IsFile(src) {
		forceExit("The input is neither an URL nor a file path, quitting! Got "+src, core.EXIT_INVALID_ARGS)
	}

	if !strings.HasSuffix(src, ".tar.xz") && !strings.HasSuffix(src, ".tar.gz") {
		forceExit("The input must be a tar.xz or tar.gz archive, quitting! Got "+src, core.EXIT_INVALID_ARGS)
	}

	log.Info("Creating %s from %s...", glog.Auto(name), glog.Auto(src))
	must(unpack(arg...))
	log.Info("Done, to enter the devbox, type:")
	log.Info(
		"sudo %s %s %s",
		glog.Bold()+core.APP_NAME+glog.Reset(),
		glog.Underline()+glog.Auto(ENTER)+glog.Reset(),
		glog.Auto(name),
	)
	return nil
}
