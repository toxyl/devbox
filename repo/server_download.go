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

func (s *Server) handleDownload(conn net.Conn, scanner *bufio.Scanner) {
	fileName, _ := ScanToken(scanner)
	fileHash, _ := ScanToken(scanner)
	user, _ := ScanToken(scanner)
	password, _ := ScanToken(scanner)

	if !s.isUser(user, password) {
		fmt.Fprintln(conn, "ERROR Access denied!")
		return
	}

	logServer.Info(
		"%s (%s) wants to download %s (%s)",
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
		glog.File(fileName),
		glog.Highlight(fileHash),
	)

	filePath := filepath.Join(s.basePath, fileName)
	fileHashServer := utils.FileToSha256(filePath)

	if fileHashServer == "0" {
		fmt.Fprintln(conn, "ERROR File not found")
		return
	}

	if fileHashServer == fileHash {
		// File is already up to date
		logServer.OK(
			"%s (%s) already has the latest version of %s",
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
			glog.File(filePath),
		)
		fmt.Fprintln(conn, "UP-TO-DATE")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(conn, "ERROR File not found")
		return
	}
	defer file.Close()

	// Seek to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Fprintln(conn, "ERROR Failed to seek file")
		return
	}

	fmt.Fprintln(conn, "FILE", fileHashServer)

	// Wait for client to respond that it's ready
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logServer.Error(
			"Failed to receive ready response from %s (%s)",
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
		)
		return
	}
	if strings.TrimSpace(response) != "READY" {
		logServer.Error(
			"Unexpected response from %s (%s): %s",
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
			glog.Highlight(response),
		)
		return
	}

	logServer.Info(
		"Sending %s to %s (%s)",
		glog.File(filePath),
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
	)
	n, err := io.Copy(conn, file)
	if err != nil {
		logServer.Error(
			"Failed to copy %s to %s (%s): %s",
			glog.File(filePath),
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
			glog.Error(err),
		)
		fmt.Fprintf(conn, "ERROR %s\n", err.Error())
		return
	}

	logServer.Success(
		"Sent %s (%s) to %s (%s)",
		glog.File(filePath),
		glog.HumanReadableBytesIEC(n),
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
	)
}
