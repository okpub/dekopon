package network

import (
	"net"

	"golang.org/x/net/websocket"
)

func DialTcp(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

func DialWeb(addr string) (net.Conn, error) {
	var (
		httpAddr  = "http://" + addr
		webAddr   = "ws://" + addr
		protocol  = ""
		conn, err = websocket.Dial(webAddr, protocol, httpAddr)
	)
	return WrapWeb(conn), err
}
