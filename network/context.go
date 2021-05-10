package network

import (
	"net"
	"time"

	"github.com/okpub/dekopon/actor"
)

//最好不要建立child
type SocketContext struct {
	actor.ActorContext
	conn net.Conn
}

func NewContext(ctx actor.ActorContext, conn net.Conn) actor.ActorContext {
	return &SocketContext{ActorContext: ctx, conn: conn}
}

func (ctx *SocketContext) Sender() actor.PID {
	return ctx.Self()
}

func (ctx *SocketContext) SetReceiveTimeout(dur time.Duration) {
	if dur > 0 {
		ctx.conn.SetReadDeadline(time.Now().Add(dur))
	} else {
		ctx.conn.SetReadDeadline(time.Time{})
	}
}

func (ctx *SocketContext) CancelReceiveTimeout() {
	ctx.SetReceiveTimeout(0)
}

func (ctx *SocketContext) Close() (err error) {
	err = ctx.conn.Close()
	return
}
