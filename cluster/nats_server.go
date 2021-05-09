package cluster

import (
	"context"
	"net"

	"github.com/skimmer/actor"
	"github.com/skimmer/conn/message"
	"github.com/skimmer/conn/packet"
	"github.com/skimmer/network"
	"github.com/skimmer/router"
)

//订阅服务
type NatsRPCServer struct {
	network.Server
	//注册其他服务器
	router.DNSRouter
	listener actor.PID
}

func NewNatServer(pid actor.PID, args ...network.ServerOption) *NatsRPCServer {
	var options = &network.ServerOptions{
		Network: network.TCP,
		MaxConn: 1000,
	}
	options.Filler(args)

	return &NatsRPCServer{Server: network.FromServer(options), listener: pid}
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
	var (
		socket = network.WithSocket(conn)
		pid    = actor.NewPID(server, actor.NewProcess(socket), actor.SetAddr(socket.Addr))
	)
	var ping = true
	var props = network.From(conn, func(ctx actor.ActorContext) {
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

	socket.RegisterHander(actor.NewSelf(pid, props))
	socket.Serve(server, conn, nil)
}

func (s *NatsRPCServer) Disconnect() {
	s.Server.Close()
}

//订阅
func (this *NatsRPCServer) Subscribe() {

}

//开启服务器
func (this *NatsRPCServer) Start() {

}

//断开服务器
func (this *NatsRPCServer) Disconnect() {

}

//推送消息(通过服务器id内部建立连接)
func (this *NatsRPCServer) PublishMessage() {

}
