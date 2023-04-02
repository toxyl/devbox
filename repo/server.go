package repo

import (
	"bufio"
	"fmt"
	"net"

	"github.com/toxyl/glog"
)

func ScanToken(scanner *bufio.Scanner) (string, bool) {
	more := scanner.Scan()
	w := scanner.Text()
	return w, more
}

type user struct {
	admin    bool
	name     string
	password string
}

type Server struct {
	basePath string
	users    []user
}

func NewServer(basePath string) *Server {
	return &Server{
		basePath: basePath,
		users:    []user{},
	}
}

func (s *Server) AddAdmin(name, password string) {
	for i, u := range s.users {
		if u.name == name {
			s.users[i].admin = true
			s.users[i].password = password
			return
		}
	}
	s.users = append(s.users, user{
		admin:    true,
		name:     name,
		password: password,
	})
}

func (s *Server) AddUser(name, password string) {
	for i, u := range s.users {
		if u.name == name {
			s.users[i].admin = false
			s.users[i].password = password
			return
		}
	}
	s.users = append(s.users, user{
		admin:    false,
		name:     name,
		password: password,
	})
}

func (s *Server) isAdmin(name, password string) bool {
	for _, u := range s.users {
		if u.name == name && u.admin && u.password == password {
			return true
		}
	}

	return false
}

func (s *Server) isUser(name, password string) bool {
	for _, u := range s.users {
		if u.name == name && u.password == password {
			return true
		}
	}
	return false
}

func (s *Server) ListenAndServe(addr string) error {
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
