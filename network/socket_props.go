package network

import (
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/utils"
)

func SetContextMiddleware(conn net.Conn) actor.ContextDecorator {
	return func(next actor.ContextDecoratorFunc) actor.ContextDecoratorFunc {
		return func(ctx actor.ActorContext) actor.ActorContext {
			return next(&SocketContext{ActorContext: ctx, conn: conn})
		}
	}
}

func fromConn(conn net.Conn, method actor.ActorFunc) *actor.Props {
	return actor.From(method).WithContextDecorator(SetContextMiddleware(conn))
}

func SpawnConn(conn net.Conn, method actor.ActorFunc, args ...SocketOption) {
	var (
		done   = utils.MakeDone()
		socket = WithSocket(conn, args...)
		pid    = actor.NewPID(actor.NewDfaultProcess(done.Done(), socket), actor.SetName(socket.Addr))
		props  = fromConn(conn, method)
	)
	defer done.Shutdown()
	socket.RegisterHander(actor.NewSelf(pid, props), props.GetDispatcher())
	socket.InvokeSystemMessage(actor.EVENT_START)
	socket.ServeConn(conn)
	socket.InvokeSystemMessage(EVENT_CLOSED)
	socket.InvokeSystemMessage(actor.EVENT_STOP)
}

func FromDial(args ...SocketOption) *actor.Props {
	return actor.From(nil).WithSpawnFunc(wrapAddrSpawner(args...))
}

func wrapAddrSpawner(args ...SocketOption) actor.SpawnFunc {
	return func(parent actor.SpawnContext, props *actor.Props, options *actor.ActorOptions) actor.PID {
		var (
			done   = utils.MakeDone()
			socket = NewSocket(args...)
			pid    = actor.NewPID(actor.NewDfaultProcess(done.Done(), socket), actor.SetName(socket.Addr))
			ctx    = actor.NewContext(parent.ChildOf(pid), props)
		)

		socket.RegisterHander(ctx, props.GetDispatcher())

		go func() {
			defer done.Shutdown()
			socket.InvokeSystemMessage(actor.EVENT_START)
			socket.Start()
			socket.InvokeSystemMessage(EVENT_CLOSED)
			socket.InvokeSystemMessage(actor.EVENT_STOP)
		}()

		return pid
	}
}
