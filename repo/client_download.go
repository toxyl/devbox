package repo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func (c *Client) DownloadFile(fileNameRemote, fileNameLocal, storagePath string) error {
	fileNameBase := filepath.Base(fileNameLocal)
	fileNameLocal = filepath.Join(storagePath, fileNameBase)
	fileHash := utils.FileToSha256(fileNameLocal)
	fmt.Fprintln(c.conn, "DOWNLOAD", filepath.Base(fileNameRemote), fileHash, c.user, c.password)
	return c.download(fileNameLocal)
}

func (c *Client) download(filePath string) error {
	// Print server response
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return err
	}
	if strings.HasPrefix(response, "ERROR") {
		return fmt.Errorf(response)
	}
	if strings.HasPrefix(response, "FILE") {
		// got a file response
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Send ready message to server
		fmt.Fprintln(c.conn, "READY")

		logClient.Info(
			"Downloading %s from %s",
			glog.File(filePath),
			glog.ConnRemote(c.conn, false),
		)
		n, err := io.Copy(file, c.conn)
		if err != nil {
			return err
		}
		logClient.Success(
			"Downloaded %s (%s) from %s",
			glog.File(filePath),
			glog.HumanReadableBytesIEC(n),
			glog.ConnRemote(c.conn, false),
		)
	}

	return nil
}
