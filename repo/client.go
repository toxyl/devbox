package repo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

type Client struct {
	conn     net.Conn
	user     string
	password string
}

func NewClient(user, password string) *Client {
	return &Client{
		user:     user,
		password: password,
	}
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) DownloadFile(fileName string) error {
	fileHash := utils.FileToSha256(fileName)
	fileName = filepath.Base(fileName)
	fmt.Fprintln(c.conn, "DOWNLOAD", fileName, fileHash)
	return c.download(fileName)
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

func (c *Client) UploadFile(fileNameSrc, fileNameDst string) error {
	f, err := os.Open(fileNameSrc)
	if err != nil {
		return err
	}
	defer f.Close()
	fileHash := utils.FileToSha256(fileNameSrc)
	fileNameSrc = filepath.Base(fileNameSrc)
	fmt.Fprintln(c.conn, "UPLOAD", fileNameDst, fileHash, c.user, utils.StringToSha256(c.password))
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

func (c *Client) Close() {
	c.conn.Close()
}
