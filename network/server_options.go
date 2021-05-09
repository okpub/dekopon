package network

import (
	"fmt"

	"github.com/skimmer/actor"
)

const (
	DefaultServerAddr = ":8999"
	DefaultMaxConn    = 1000
)

func NewServer(args ...ServerOption) Server {
	var options = ServerOptions{
		MaxConn: DefaultMaxConn,
		Addr:    DefaultServerAddr,
		Network: TCP,
	}
	options.Filler(args)

	switch options.Network {
	case TCP:
		return &TcpServer{ServerOptions: options, TaskDone: actor.MakeDone()}
	case WEB:
		return &WebServer{ServerOptions: options, TaskDone: actor.MakeDone()}
	default:
		panic(fmt.Errorf("the network untype: %s", options.Network))
	}
}

type ServerOptions struct {
	Addr     string
	Network  string
	MaxConn  int
	ReadSize int
}

func (options ServerOptions) Options() ServerOptions { return options }

func (options *ServerOptions) Filler(args []ServerOption) *ServerOptions {
	for _, f := range args {
		f(options)
	}
	return options
}

//可选参数
type ServerOption func(*ServerOptions)

func SetAddr(addr string) ServerOption {
	return func(p *ServerOptions) {
		p.Addr = addr
	}
}

func SetNetwork(addrType string) ServerOption {
	return func(p *ServerOptions) {
		p.Network = addrType
	}
}

func SetMaxConn(n int) ServerOption {
	return func(p *ServerOptions) {
		p.MaxConn = n
	}
}
