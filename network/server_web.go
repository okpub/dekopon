package network

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/okpub/dekopon/utils"
	"golang.org/x/net/websocket"
)

type WebServer struct {
	utils.TaskDone
	*ServerOptions
}

func (server *WebServer) Close() (err error) {
	server.Shutdown()
	return
}

func (server *WebServer) ListenAndServe(ctx context.Context, handler Handler) (err error) {
	var (
		child, cancel = context.WithCancel(ctx)
		conns         ConnSet
		handleWebConn = func() websocket.Handler {
			return func(conn *websocket.Conn) {
				if conns.SetConnMax(conn, server.MaxConn) {
					handler(child, WrapWeb(conn))
				} else {
					fmt.Println("WARNING: max full conn")
				}
			}
		}
		ln = &http.Server{Addr: server.Addr, Handler: handleWebConn()}
	)

	go func() {
		defer ln.Shutdown(child)
		select {
		case <-server.Done():
			cancel()
		case <-child.Done():
			//todo
		}
	}()

	func() {
		defer server.Close()
		defer conns.CloseAll()
		err = ln.ListenAndServe()
		fmt.Println("EXIT: close web server#", err)
	}()

	return
}

//web conn
type WebConn struct {
	*websocket.Conn
}

func WrapWeb(conn *websocket.Conn) net.Conn {
	return &WebConn{Conn: conn}
}

func (conn *WebConn) Write(body []byte) (n int, err error) {
	if err = websocket.Message.Send(conn.Conn, body); err == nil {
		n = len(body)
	}
	return
}
