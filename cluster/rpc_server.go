package cluster

import (
	"context"
	"fmt"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/conn/message"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/network"
)

//应答服务器
type RPCServer struct {
	network.Server
	PIDSet
}

func NewRPCServer(pid PIDSet, args ...network.ServerOption) *RPCServer {
	var options = &network.ServerOptions{
		Network: network.TCP,
		MaxConn: 1000,
	}
	options.Filler(args)
	return &RPCServer{Server: network.FromServer(options), PIDSet: pid}
}

func (s *RPCServer) Start(ctx context.Context) {
	go s.Serve(ctx)
}

func (s *RPCServer) Serve(ctx context.Context) (err error) {
	err = s.Server.ListenAndServe(ctx, s.handleConn)
	return
}

func (s *RPCServer) Disconnect() {
	s.Server.Close()
}

func (s *RPCServer) handleConn(server context.Context, conn net.Conn) {
	var (
		socket = network.WithSocket(conn)
		pid    = actor.NewPID(server, actor.NewProcess(socket), actor.SetAddr(socket.Addr))
	)
	var props = network.From(conn, func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *packet.Packet:
			s.processExchange(ctx, message.Request(event))
		}
	})

	socket.RegisterHander(actor.NewSelf(pid, props))
	socket.ServeConn(server, conn)
}

//通过serverName获取不同通道的pid,将多个或者一个功能集合
func (s *RPCServer) processExchange(ctx actor.ActorContext, req *rpc.Request) {
	if pid, ok := s.PIDSet[req.ServerName]; ok {
		if res, err := pid.Call(req); err == nil {
			ctx.Respond(res)
		} else {
			ctx.Stop(ctx.Self())
		}
	} else {
		fmt.Println("ERROR: can't find ServerName:", req.ServerName)
		ctx.Stop(ctx.Self())
	}
}
