package network

import (
	"net"

	"github.com/okpub/dekopon/actor"
)

func SetContextMiddleware(conn net.Conn) actor.ContextDecorator {
	return func(next actor.ContextDecoratorFunc) actor.ContextDecoratorFunc {
		return func(ctx actor.ActorContext) actor.ActorContext {
			return next(&SocketContext{ActorContext: ctx, conn: conn})
		}
	}
}

func From(conn net.Conn, method actor.ActorFunc) *actor.Props {
	return actor.From(method).WithContextDecorator(SetContextMiddleware(conn))
}

func FromDial(args ...SocketOption) *actor.Props {
	return actor.From(nil).WithSpawnFunc(wrapAddrSpawner(args...))
}

func wrapAddrSpawner(args ...SocketOption) actor.SpawnFunc {
	return func(parent actor.SpawnContext, props *actor.Props, options *actor.PIDOptions) actor.PID {
		var (
			child, cancel = props.WithCancel(parent.Background())
			socket        = NewSocket(args...)
			pid           = actor.NewPID(child, actor.NewProcess(socket), actor.SetAddr(socket.Addr))
			ctx           = actor.NewContext(parent.ChildOf(pid), props)
		)

		socket.RegisterHander(ctx)
		//props可以不融入conn (只是为了设置读取超时)
		go func() {
			defer cancel()
			socket.Start(child)
		}()

		return pid
	}
}
