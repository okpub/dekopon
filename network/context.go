package network

import (
	"net"
	"time"

	"github.com/skimmer/actor"
)

//extends ActorContext
type SocketContext struct {
	actor.ActorContext
	conn net.Conn
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
