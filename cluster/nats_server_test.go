package cluster

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/skimmer/actor"
	"github.com/skimmer/bean/cmd"
	"github.com/skimmer/bean/message/common"
	"github.com/skimmer/bean/message/login"
	"github.com/skimmer/conn/message"
	"github.com/skimmer/conn/packet"
	"github.com/skimmer/network"
)

func test_client(uid int32) {
	var (
		socket = network.NewSocket(network.SetDialAddr("localhost:9090"))
		pid    = actor.NewPID(context.Background(), actor.NewProcess(socket), actor.SetAddr(socket.Addr))
	)
	fmt.Println("建立socket:", pid)
	go func() {
		//1 先连接
		var conn, err = socket.Connect()

		//2 连接后建立context(无论连接是否成功，都建立)
		var props = network.From(conn, func(ctx actor.ActorContext) {
			switch event := ctx.Message().(type) {
			case *packet.Packet:
				var msg = message.UnPack(event)
				fmt.Println(uid, "客户端收到消息:", msg.Header)
			}
		})

		//3 包装成context
		socket.RegisterHander(actor.NewSelf(pid, props))
		socket.Serve(context.Background(), conn, err)
	}()

	var req = message.Pack(cmd.CMD_LOGIN, message.SetMessageData(&login.LoginReq{UserID: uid}))
	pid.Send(req)
}

type TestActor struct{}

func NewTestActor() actor.Actor {
	return &TestActor{}
}

func (a *TestActor) Received(ctx actor.ActorContext) {
	switch event := ctx.Message().(type) {
	case *actor.Started:
		a.start(ctx)
	case *common.UserMessage:
		switch event.Header.Cmd {
		case cmd.CMD_LOGIN:
			var data = &login.LoginReq{}
			message.GetMessage(event, data)
			fmt.Println("处理登陆:", event.Header.Cmd, data)
			var p1 = message.SetMessageData(&login.LoginResp{ServerID: 110, RoomID: 2001})
			ctx.Respond(message.Pack(cmd.CMD_LOGIN, p1))
		}
	}
}

func (*TestActor) start(ctx actor.ActorContext) {
	var svr = NewNatServer(ctx.Self(), network.SetAddr(":9090"))
	svr.Start(ctx.Background())
}

func TestInit(t *testing.T) {
	var root, _ = context.WithTimeout(context.Background(), time.Second*3)
	var stage = actor.WithSystem(root)

	stage.ActorOf(actor.FromProducer(NewTestActor))

	time.Sleep(time.Millisecond * 10)

	var i int32
	for i = 100; i < 110; i++ {
		test_client(i)
	}

	stage.Wait()
	fmt.Println("结束")
}
