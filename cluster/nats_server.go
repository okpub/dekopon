package cluster

import (
	"context"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/conn/message"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/network"
	"github.com/okpub/dekopon/router"
)

//订阅服务
type NatsRPCServer struct {
	network.Server
	//注册其他服务器
	router.DNSRouter
	listener actor.PID
}

func NewNatServer(pid actor.PID, args ...network.ServerOption) *NatsRPCServer {
	return &NatsRPCServer{Server: network.NewServer(args...), listener: pid}
}

func (s *NatsRPCServer) Start(ctx context.Context) {
	go s.Serve(ctx)
}

func (s *NatsRPCServer) Serve(ctx context.Context) (err error) {
	err = s.Server.ListenAndServe(ctx, s.handleConn)
	return
}

//无状态控制
func (s *NatsRPCServer) handleConn(server context.Context, conn net.Conn) {
	var ping = false
	network.SpawnConn(server, conn, func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *network.EventOpen:
			ctx.SetReceiveTimeout(network.PingTime)
		case *packet.Packet:
			ping = true
			ctx.SetReceiveTimeout(network.PingTime)
			ctx.Request(s.listener, message.UnPack(event))
		case *network.TempErr:
			ctx.SetReceiveTimeout(network.PingTime)
			if ping {
				ping = false
				ctx.Request(s.listener, event)
			} else {
				ctx.Stop(ctx.Self())
			}
		}
	})
}

func (s *NatsRPCServer) Disconnect() {
	s.Server.Close()
}
