package network

import "context"

//服务器选项
type ServerOptions struct {
	context.Context
	Addr    string
	Network string
	MaxConn int
}

func (options ServerOptions) Options() *ServerOptions { return &options }

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

func SetServerBackground(ctx context.Context) ServerOption {
	return func(p *ServerOptions) {
		p.Context = ctx
	}
}
