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

var wgRoot actor.WaitGroup

func handlerConn(parent actor.NodePart, handler actor.PID) Handler {
	return func(server context.Context, conn net.Conn) {
		var ping = false
		SpawnConn(server, conn, func(ctx actor.ActorContext) {
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
			var pid = FromServer(ctx.Background(), handlerConn(stage, router[3]), SetAddr(":9001"), SetNetwork(WEB))
			ctx.Watch(pid)

			time.Sleep(time.Millisecond * 10)
			dial_client(stage, "localhost:9001", WEB)

			time.Sleep(time.Second * 1)
			pid.Close()
		case *actor.Terminated:
			fmt.Println("被迫营业:", event.GetWho())
		case *actor.Stopped:
			//可以等待服务器pid
		}
	}))

	var tcpServer = NewServer(SetAddr(":9003"))
	wgRoot.Wrap(func() {
		tcpServer.ListenAndServe(ccc, handlerConn(stage, router[3]))
	})
	time.Sleep(time.Millisecond * 10)
	dial_client(stage, "localhost:9003", TCP)

	wgRoot.Wait()

	stage.Wait()
	fmt.Println("结束")
}
