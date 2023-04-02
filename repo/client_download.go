package repo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

type ProgressWriter struct {
	Total          int64
	ReportInterval int64
	Callback       func(total int64)
	w              io.Writer
	lastReportTime time.Time
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.w.Write(p)
	pw.Total += int64(n)
	if int64(time.Since(pw.lastReportTime).Seconds()) > pw.ReportInterval && pw.Callback != nil {
		pw.Callback(pw.Total)
		pw.lastReportTime = time.Now()
		if fw, ok := pw.w.(*os.File); ok {
			if err := fw.Sync(); err != nil {
				return n, err
			}
		}
	}
	return
}

func (pw *ProgressWriter) Wrap(w io.Writer) io.Writer {
	pw.w = w
	return pw
}

// DownloadFile tries to download the file. If it returns `true, nil` then a new file has been downloaded.
// Otherwise there was either an error or the file hasn't changed.
func (c *Client) DownloadFile(fileNameRemote, fileNameLocal, storagePath string) (bool, error) {
	fileNameBase := filepath.Base(fileNameLocal)
	fileNameLocal = filepath.Join(storagePath, fileNameBase)
	fileHash := utils.FileToSha256(fileNameLocal)
	fmt.Fprintln(c.conn, "DOWNLOAD", filepath.Base(fileNameRemote), fileHash, c.user, c.password)
	return c.download(fileNameLocal)
}

func (c *Client) download(filePath string) (bool, error) {
	// Print server response
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(response, "ERROR") {
		return false, fmt.Errorf(response)
	}
	if strings.HasPrefix(response, "FILE") {
		// got a file response
		file, err := os.Create(filePath)
		if err != nil {
			return false, err
		}
		defer file.Close()

		// Send ready message to server
		fmt.Fprintln(c.conn, "READY")

		logClient.Info(
			"Downloading %s from %s...",
			glog.File(filePath),
			glog.ConnRemote(c.conn, false),
		)

		c.progressWriter = &ProgressWriter{
			Total:          0,
			ReportInterval: 10, // report progress every 10s
			Callback: func(total int64) {
				logClient.Info(
					"Downloading %s from %s: %s...",
					glog.File(filePath),
					glog.ConnRemote(c.conn, false),
					glog.HumanReadableBytesIEC(total),
				)
			},
		}

		n, err := io.Copy(c.progressWriter.Wrap(file), c.conn)
		if err != nil {
			return false, err
		}
		logClient.Success(
			"Downloaded  %s from %s: %s total",
			glog.File(filePath),
			glog.ConnRemote(c.conn, false),
			glog.HumanReadableBytesIEC(n),
		)
		return true, nil
	}

	return false, nil
}
