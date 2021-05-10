package network

import (
	"context"
	"fmt"

	"github.com/okpub/dekopon/actor"
)

const (
	DefaultServerAddr = ":8999"
	DefaultMaxConn    = 1000
)

func NewServerOptions(args ...ServerOption) *ServerOptions {
	var options = ServerOptions{
		MaxConn: DefaultMaxConn,
		Addr:    DefaultServerAddr,
		Network: TCP,
	}
	return options.Filler(args)
}

func NewServer(args ...ServerOption) Server {
	return WithServer(NewServerOptions(args...))
}

func WithServer(options *ServerOptions) Server {
	switch options.Network {
	case TCP:
		return &TcpServer{ServerOptions: options, TaskDone: actor.MakeDone()}
	case WEB:
		return &WebServer{ServerOptions: options, TaskDone: actor.MakeDone()}
	default:
		panic(fmt.Errorf("can't open server untype: %s", options.Network))
	}
}

func FromServer(ctx context.Context, handler Handler, args ...ServerOption) actor.PID {
	var (
		options       = NewServerOptions(args...)
		child, cancel = context.WithCancel(ctx)
		ref           = NewServerProcess(options, handler)
		pid           = actor.NewPID(child, ref, actor.SetAddr(options.Addr))
	)

	go func() {
		defer cancel()
		ref.PostSystemMessage(pid, actor.EVENT_START)
		ref.PostSystemMessage(pid, actor.EVENT_STOP)
	}()

	return pid
}
