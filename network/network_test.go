package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/login"
	"github.com/okpub/dekopon/conn/message"
	"github.com/okpub/dekopon/conn/packet"
)

func handlerConn(handler actor.PID) Handler {
	return func(server context.Context, conn net.Conn) {
		var ping = false
		SpawnConn(conn, func(ctx actor.ActorContext) {
			switch event := ctx.Message().(type) {
			case *EventOpen:
				ctx.SetReceiveTimeout(time.Second * 1)
			case *packet.Packet:
				ping = true
				ctx.SetReceiveTimeout(PingTime)
				var msg = message.UnPack(event)
				fmt.Println("服务端收到消息:", msg.Header)
				ctx.RequestWithCustomSender(handler, event, ctx.Sender())
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
	}
}

func dial_client(parent actor.SpawnContext, addr, kind string) {
	var (
		p = FromDial(SetDialAddr(addr), SetDialNetwork(kind))
	)

	p.WithFunc(func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *packet.Packet:
			var msg = message.UnPack(event)
			fmt.Println("客户端收到消息:", msg.Header)
		case *DialError:
		case *EventOpen:
		case *EventClose:
		}
	})

	var pid = parent.ActorOf(p)
	if kind == TCP {
		var data = message.Pack(101, message.SetMessageType(102), message.SetMessageData(&login.LoginReq{Pwd: "密码"}))
		pid.Send(data)
	}
}

var (
	router = map[int]actor.PID{}
)

func register_serve(parent actor.SpawnContext, id int, name string) {
	router[id] = parent.ActorOf(actor.From(func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *actor.Started:
		case *actor.Stopped:
			fmt.Println(name, "退出")
		case *packet.Packet:
			var msg = message.UnPack(event)

			fmt.Println(name, "路由收到消息:", msg)
			ctx.Respond(message.Pack(999))
		}
	}))
}

func TestInit(t *testing.T) {
	var (
		ccc, cancel = context.WithTimeout(context.Background(), time.Second*3)
		stage       = actor.WithSystem(ccc)
	)
	defer cancel()

	register_serve(stage, 1, "大厅")
	register_serve(stage, 2, "房间")
	register_serve(stage, 3, "登陆")

	stage.ActorOf(actor.From(func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *actor.Started:
			var pid = FromServer(
				handlerConn(router[3]),
				SetAddr(":9001"),
				SetNetwork(WEB))

			ctx.Watch(pid)

			time.Sleep(time.Millisecond * 10)
			dial_client(stage, "localhost:9001", WEB)

			time.Sleep(time.Second * 1)
			pid.GracefulStop()
			fmt.Println("关闭服务器")
		case *actor.Terminated:
			fmt.Println("被迫营业:", event.GetWho())
		case *actor.Stopped:
			//可以等待服务器pid
		}
	}))

	FromServer(handlerConn(router[3]), SetAddr(":9003"), SetServerBackground(ccc))

	time.Sleep(time.Millisecond * 10)
	dial_client(stage, "localhost:9003", TCP)

	fmt.Println("结束1")
	stage.Wait()
	fmt.Println("结束2")
}
