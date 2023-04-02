package repo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func (s *Server) handleUpload(conn net.Conn, scanner *bufio.Scanner) {
	// Receive file name and hash from client
	fileName, _ := ScanToken(scanner)
	fileHash, _ := ScanToken(scanner)
	user, _ := ScanToken(scanner)
	password, _ := ScanToken(scanner)

	if !s.isAdmin(user, password) {
		fmt.Fprintln(conn, "ERROR Access denied!")
		return
	}

	logServer.Info(
		"%s (%s) wants to upload %s (%s)",
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
		glog.File(fileName),
		glog.Highlight(fileHash),
	)

	filePath := filepath.Join(s.basePath, fileName)
	fileHashServer := utils.FileToSha256(filePath)

	if fileHashServer == fileHash {
		// File is already up to date
		logServer.OK(
			"We already have the version of %s that %s (%s) tries to upload",
			glog.File(filePath),
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
		)
		fmt.Fprintln(conn, "UP-TO-DATE")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		// create dirs and file
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			fmt.Fprintln(conn, "ERROR Could not create dir:", err.Error())
			return
		}
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Fprintln(conn, "ERROR Could not create file:", err.Error())
			return
		}
	}
	defer file.Close()
	// Seek to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Fprintln(conn, "ERROR Failed to seek file")
		return
	}

	// Signal the client that we are ready to receive
	fmt.Fprintln(conn, "READY")

	logServer.Info(
		"Receiving %s from %s (%s)",
		glog.File(filePath),
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
	)
	n, err := io.Copy(file, conn)
	if err != nil {
		logServer.Error(
			"Failed to copy %s from %s (%s): %s",
			glog.File(filePath),
			glog.Highlight(user),
			glog.ConnRemote(conn, false),
			glog.Error(err),
		)
		fmt.Fprintf(conn, "ERROR %s\n", err.Error())
		return
	}

	logServer.Success(
		"Received %s (%s) from %s (%s)",
		glog.File(filePath),
		glog.HumanReadableBytesIEC(n),
		glog.Highlight(user),
		glog.ConnRemote(conn, false),
	)
}
