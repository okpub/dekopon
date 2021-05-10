package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/mailbox"
)

type TcpServer struct {
	actor.TaskDone
	*ServerOptions
}

/*
 * 1 这里之所以新建一context是因为服务器退出，他所有依赖的子context退出
 * 2 如果服务器退出，并不意味此服务关闭，可以循环调用来重启
 */
func (server *TcpServer) ListenAndServe(ctx context.Context, handler Handler) (err error) {
	var (
		ln net.Listener
	)

	if ln, err = net.Listen(TCP, server.Addr); mailbox.Fail(err) {
		fmt.Println("ERROR: close tcp server#", err)
		return
	}

	var (
		wg            actor.WaitGroup
		conns         ConnSet
		child, cancel = context.WithCancel(ctx)
		handlerConn   = func(conn net.Conn) func() {
			return func() {
				handler(child, conn)
				conns.RemoveConn(conn)
			}
		}
	)

	go func() {
		defer ln.Close()
		select {
		case <-server.Done():
		case <-child.Done():
		}
	}()

	func() {
		var conn net.Conn
		defer cancel()
		defer conns.CloseAll()
		for {
			if conn, err = ln.Accept(); err == nil {
				if conns.SetConnMax(conn, server.MaxConn) {
					wg.Wrap(handlerConn(conn))
				} else {
					conn.Close()
					fmt.Println("WARNING: max full conn")
				}
			} else {
				if Temporary(err) {
					time.Sleep(TemporaryInterval)
				} else {
					break
				}
			}
		}
	}()

	wg.Wait()

	fmt.Println("EXIT: close tcp server#", err)
	return
}
