package repo

import (
	"net"
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

func (c *Client) Close() {
	c.conn.Close()
}
