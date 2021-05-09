package network

import (
	"net"
	"time"
)

type Connecter struct {
	conn net.Conn
	err  error
}

func New(conn net.Conn) net.Conn {
	return &Connecter{conn: conn}
}

func NewErr(conn net.Conn, err error) *Connecter {
	return &Connecter{conn: conn, err: err}
}

func (c *Connecter) Err() error {
	return c.err
}

func (c *Connecter) Read(b []byte) (n int, err error) {
	if err = c.err; err == nil {
		n, err = c.conn.Read(b)
	}
	return
}

func (c *Connecter) Write(b []byte) (n int, err error) {
	if err = c.err; err == nil {
		n, err = c.conn.Write(b)
	}
	return
}

func (c *Connecter) Close() (err error) {
	if err = c.err; err == nil {
		err = c.conn.Close()
	}
	return
}

func (c *Connecter) LocalAddr() net.Addr {
	if err := c.err; err == nil {
		return c.conn.LocalAddr()
	}
	return nil
}

func (c *Connecter) RemoteAddr() net.Addr {
	if err := c.err; err == nil {
		return c.conn.RemoteAddr()
	}
	return nil
}

func (c *Connecter) SetDeadline(t time.Time) (err error) {
	if err = c.err; err == nil {
		err = c.conn.SetDeadline(t)
	}
	return
}

func (c *Connecter) SetReadDeadline(t time.Time) (err error) {
	if err = c.err; err == nil {
		err = c.conn.SetReadDeadline(t)
	}
	return
}

func (c *Connecter) SetWriteDeadline(t time.Time) (err error) {
	if err = c.err; err == nil {
		err = c.conn.SetWriteDeadline(t)
	}
	return
}
