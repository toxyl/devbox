package repo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func (c *Client) UploadFile(fileNameSrc, fileNameDst string) error {
	f, err := os.Open(fileNameSrc)
	if err != nil {
		return err
	}
	f.Close()
	fileHash := utils.FileToSha256(fileNameSrc)
	fmt.Fprintln(c.conn, "UPLOAD", fileNameDst, fileHash, c.user, c.password)
	return c.upload(fileNameSrc)
}

func (c *Client) upload(filePath string) error {
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return err
	}
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "ERROR") {
		return fmt.Errorf(response)
	}
	if response == "UP-TO-DATE" {
		logClient.Success(
			"No need to upload %s, %s already has this version",
			glog.File(filePath),
			glog.ConnRemote(c.conn, false),
		)
		return nil
	}
	if response == "READY" {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		logClient.Info(
			"Uploading %s to %s",
			glog.File(filePath),
			glog.ConnRemote(c.conn, false),
		)
		n, err := io.Copy(c.conn, file)
		if err != nil {
			return err
		}
		logClient.Success(
			"Uploaded %s (%s) to %s",
			glog.File(filePath),
			glog.HumanReadableBytesIEC(n),
			glog.ConnRemote(c.conn, false),
		)
	}

	return nil
}
