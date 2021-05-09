package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/skimmer/actor"
	"github.com/skimmer/bean/message/login"
	"github.com/skimmer/conn/message"
	"github.com/skimmer/conn/packet"
)

var wgRoot actor.WaitGroup

func handlerConn(handler actor.PID) Handler {
	return func(server context.Context, conn net.Conn) {
		var (
			options = NewOptions()
			socket  = WithSocket(options, conn)
			pid     = actor.NewPID(server, actor.NewProcess(socket), actor.SetAddr(socket.Addr))
		)
		var ping = false
		var props = From(conn, func(ctx actor.ActorContext) {
			switch event := ctx.Message().(type) {
			case *EventOpen:
				ctx.SetReceiveTimeout(time.Second * 1)
			case *packet.Packet:
				ping = true
				ctx.SetReceiveTimeout(PingTime)
				var msg = message.UnPack(event)
				fmt.Println("服务端收到消息:", msg.Header)
			case *TempErr:
				ctx.CancelReceiveTimeout()
				if ping {
					//todo
					fmt.Println("临时错误:有响应")
				} else {
					fmt.Println("临时错误:关闭")
					ctx.Stop(ctx.Self())
				}
			}
		})

		var ctx = actor.NewContext(stage.ChildOf(pid), props)

		socket.RegisterHander(ctx)

		socket.Serve(server, conn, nil)

	}
}

func dial_client(addr, kind string) {
	var options = NewOptions()
	options.WithAddr(addr).WithNetwork(kind)

	var (
		socket = NewSocket(options) //不能对外公开
		pid    = actor.NewPID(context.Background(), actor.NewProcess(socket), actor.SetAddr(socket.Addr))
	)
	fmt.Println("建立socket:", pid)
	go func() {
		//1 先连接
		var conn, err = socket.Connect()

		//2 连接后建立context(无论连接是否成功，都建立)
		var props = From(conn, func(ctx actor.ActorContext) {
			switch event := ctx.Message().(type) {
			case *packet.Packet:
				var msg = message.UnPack(event)
				fmt.Println("客户端收到消息:", msg.Header)
			case *DialError:

			case *EventOpen:

			case *EventClose:

			}
		})

		//3 包装成context
		var ctx = actor.NewContext(stage.ChildOf(pid), props)

		socket.RegisterHander(ctx)
		socket.Serve(context.Background(), conn, err)
	}()

	if kind == TCP {
		var data = message.Pack(101, message.SetType(102), message.SetMessage(&login.LoginReq{Pwd: "密码"}))
		pid.Send(data)
	}
}

var (
	router = map[int]actor.PID{}
	stage  = actor.NewSystem()
)

func register_serve(id int, name string) {
	router[id] = stage.ActorOf(actor.From(func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *actor.Started:
		case *actor.Stopped:
			fmt.Println(name, "退出")
		case actor.Request:
			var p, ok = event.Message().(*packet.Packet)
			if ok {
				var msg = message.UnPack(p)
				fmt.Println("收到同步消息:", msg)
			}
		}
	}))
}

func TestInit(t *testing.T) {
	var (
		ctx, _ = context.WithTimeout(context.Background(), time.Second*3)
	)

	go func() {
		actor.Wait(ctx)
		stage.Shutdown()
	}()

	register_serve(1, "大厅")
	register_serve(2, "房间")
	register_serve(3, "登陆")

	var tcpServer = NewServer(SetAddr(":9003"))
	wgRoot.Wrap(func() {
		tcpServer.ListenAndServe(ctx, handlerConn(router[3]))
	})
	time.Sleep(time.Millisecond * 10)
	dial_client("localhost:9003", TCP)

	var clientServer = NewServer(SetAddr(":9001"), SetNetwork(WEB))
	wgRoot.Wrap(func() {
		clientServer.ListenAndServe(ctx, handlerConn(router[3]))
	})
	time.Sleep(time.Millisecond * 10)
	dial_client("localhost:9001", WEB)
	wgRoot.Wait()

	stage.Wait()
	fmt.Println("结束")
}
