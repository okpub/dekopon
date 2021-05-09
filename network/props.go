package network

import (
	"net"

	"github.com/skimmer/actor"
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
