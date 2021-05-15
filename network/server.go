package network

import (
	"context"
	"fmt"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/utils"
)

const (
	defaultServerAddr = ":8999"
	defaultMaxConn    = 1000
)

func NewServer(args ...ServerOption) Server {
	var options = ServerOptions{
		Context: context.Background(),
		MaxConn: defaultMaxConn,
		Addr:    defaultServerAddr,
		Network: TCP,
	}
	return WithServer(options.Filler(args))
}

func WithServer(options *ServerOptions) Server {
	switch options.Network {
	case TCP:
		return &TcpServer{ServerOptions: options, TaskDone: utils.MakeDone()}
	case WEB:
		return &WebServer{ServerOptions: options, TaskDone: utils.MakeDone()}
	default:
		panic(fmt.Errorf("can't open server untype: %s", options.Network))
	}
}

func FromServer(handler Handler, args ...ServerOption) actor.PID {
	var (
		done   = utils.MakeDone()
		server = NewServer(args...)
		ref    = NewServerProcess(done, server)
		pid    = actor.NewPID(ref, actor.SetName(server.Options().Addr))
	)

	go func() {
		defer done.Shutdown()
		ref.PostSystemMessage(pid, actor.EVENT_START)
		server.ListenAndServe(server.Options().Context, handler)
		ref.PostSystemMessage(pid, actor.EVENT_STOP)
	}()

	return pid
}
