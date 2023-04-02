package repo

import (
	"bufio"
	"fmt"
	"net"
	"path/filepath"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func (s *Server) handleCheck(conn net.Conn, scanner *bufio.Scanner) {
	fileName, _ := ScanToken(scanner)
	fileHash, _ := ScanToken(scanner)
	user, _ := ScanToken(scanner)
	password, _ := ScanToken(scanner)

	if !s.isUser(user, password) {
		fmt.Fprintln(conn, "ERROR Access denied!")
		return
	}

	filePath := filepath.Join(s.basePath, fileName)
	fileHashServer := utils.FileToSha256(filePath)

	logServer.Info("Checking %s (%s): %s", glog.File(filePath), glog.Highlight(fileHashServer), glog.Highlight(fileHash))

	if fileHashServer == "0" {
		fmt.Fprintln(conn, "ERROR File not found")
		return
	}

	if fileHashServer == fileHash {
		fmt.Fprintln(conn, "UP-TO-DATE")
		return
	}

	fmt.Fprintln(conn, "OUTDATED")
}
