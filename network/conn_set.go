package network

import (
	"net"
	"sync"
)

type ConnSet struct {
	mu    sync.Mutex
	conns map[net.Conn]struct{}
}

/*
* max>-1为有限添加
 */
func (conns *ConnSet) SetConnMax(conn net.Conn, max int) (ok bool) {
	conns.mu.Lock()
	if conns.conns == nil {
		conns.conns = make(map[net.Conn]struct{})
	}
	if ok = 0 > max || max > len(conns.conns); ok {
		conns.conns[conn] = struct{}{}
	}
	conns.mu.Unlock()
	return
}

/*
* 无上限添加
 */
func (conns *ConnSet) SetConn(conn net.Conn) bool {
	return conns.SetConnMax(conn, LimitConnsNone)
}

func (conns *ConnSet) RemoveConn(conn net.Conn) {
	conns.mu.Lock()
	delete(conns.conns, conn)
	conns.mu.Unlock()
	conn.Close()
}

func (conns *ConnSet) Len() (n int) {
	conns.mu.Lock()
	n = len(conns.conns)
	conns.mu.Unlock()
	return
}

func (conns *ConnSet) CloseAll() {
	var list []net.Conn

	conns.mu.Lock()
	for conn := range conns.conns {
		list = append(list, conn)
	}
	conns.conns = nil
	conns.mu.Unlock()

	for _, conn := range list {
		conn.Close()
	}
}
