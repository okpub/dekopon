package network

import (
	"fmt"
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

func Dial(options *SocketOptions) (conn net.Conn, err error) {
	switch options.Network {
	case WEB:
		conn, err = DialWeb(options.Addr)
	case TCP:
		conn, err = DialTcp(options.Addr)
	default:
		err = fmt.Errorf("can't dial network %s", options.Network)
	}
	return
}
