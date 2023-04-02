package repo

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/toxyl/devbox/utils"
	"github.com/toxyl/glog"
)

func ScanToken(scanner *bufio.Scanner) (string, bool) {
	more := scanner.Scan()
	w := scanner.Text()
	return w, more
}

type Server struct {
	basePath string
	user     string
	password string
}

func NewServer(user, password string) *Server {
	return &Server{
		user:     user,
		password: password,
	}
}

func (s *Server) ListenAndServe(addr string, basePath string) error {
	s.basePath = basePath
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	logServer.OK("Server listening on %s", glog.Addr(addr, false))

	for {
		conn, err := l.Accept()
		if err != nil {
			logServer.Error("Failed to accept connection: %s", glog.Error(err))
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleDownload(conn net.Conn, scanner *bufio.Scanner) {
	fileName, _ := ScanToken(scanner)
	fileHash, _ := ScanToken(scanner)

	logServer.Info(
		"%s wants to download %s (%s)",
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
			"%s already has the latest version of %s",
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
			"Failed to receive ready response from %s",
			glog.ConnRemote(conn, false),
		)
		return
	}
	if strings.TrimSpace(response) != "READY" {
		logServer.Error(
			"Unexpected response from %s: %s",
			glog.ConnRemote(conn, false),
			glog.Highlight(response),
		)
		return
	}

	logServer.Info(
		"Sending %s to %s",
		glog.File(filePath),
		glog.ConnRemote(conn, false),
	)
	n, err := io.Copy(conn, file)
	if err != nil {
		logServer.Error(
			"Failed to copy %s to %s: %s",
			glog.File(filePath),
			glog.ConnRemote(conn, false),
			glog.Error(err),
		)
		fmt.Fprintf(conn, "ERROR %s\n", err.Error())
		return
	}

	logServer.Success(
		"Sent %s (%s) to %s",
		glog.File(filePath),
		glog.HumanReadableBytesIEC(n),
		glog.ConnRemote(conn, false),
	)
}

func (s *Server) handleUpload(conn net.Conn, scanner *bufio.Scanner) {
	// Receive file name and hash from client
	fileName, _ := ScanToken(scanner)
	fileHash, _ := ScanToken(scanner)
	user, _ := ScanToken(scanner)
	passwordHash, _ := ScanToken(scanner)

	if user != s.user || passwordHash != utils.StringToSha256(s.password) {
		fmt.Fprintln(conn, "ERROR Access denied!")
		return
	}

	logServer.Info(
		"%s wants to upload %s (%s)",
		glog.ConnRemote(conn, false),
		glog.File(fileName),
		glog.Highlight(fileHash),
	)

	// Check if file exists and is within basePath
	filePath := filepath.Join(s.basePath, fileName)
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

	// Compute hash of file
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Fprintln(conn, "ERROR Failed to compute hash")
		return
	}
	fileHashServer := hex.EncodeToString(hash.Sum(nil))

	if fileHashServer == fileHash {
		// File is already up to date
		logServer.OK(
			"We already have the version of %s that %s tries to upload",
			glog.File(filePath),
			glog.ConnRemote(conn, false),
		)
		fmt.Fprintln(conn, "UP-TO-DATE")
		return
	}

	// Seek to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Fprintln(conn, "ERROR Failed to seek file")
		return
	}

	// Signal the client that we are ready to receive
	fmt.Fprintln(conn, "READY")

	logServer.Info(
		"Receiving %s from %s",
		glog.File(filePath),
		glog.ConnRemote(conn, false),
	)
	n, err := io.Copy(file, conn)
	if err != nil {
		logServer.Error(
			"Failed to copy %s from %s: %s",
			glog.File(filePath),
			glog.ConnRemote(conn, false),
			glog.Error(err),
		)
		fmt.Fprintf(conn, "ERROR %s\n", err.Error())
		return
	}

	logServer.Success(
		"Received %s (%s) from %s",
		glog.File(filePath),
		glog.HumanReadableBytesIEC(n),
		glog.ConnRemote(conn, false),
	)
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	action, _ := ScanToken(scanner)
	switch action {
	case "DOWNLOAD":
		s.handleDownload(conn, scanner)
	case "UPLOAD":
		s.handleUpload(conn, scanner)
	default:
		fmt.Fprintln(conn, "ERROR Unknown command")
		return
	}
}
