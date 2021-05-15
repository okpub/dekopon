package cluster

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/cmd"
	"github.com/okpub/dekopon/bean/message/common"
	"github.com/okpub/dekopon/bean/message/login"
	"github.com/okpub/dekopon/conn/message"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/network"
)

func test_client(parent actor.SpawnContext, uid int32) {
	var p = network.FromDial(network.SetDialAddr("localhost:9090"))
	var pid = parent.ActorOf(p.WithFunc(func(ctx actor.ActorContext) {
		switch event := ctx.Message().(type) {
		case *packet.Packet:
			var msg = message.UnPack(event)
			fmt.Println(uid, "客户端收到消息:", msg.Header)
		}
	}))

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
	svr.Start(context.Background())
}

func TestInit2(t *testing.T) {
	var (
		root, cancel = context.WithTimeout(context.Background(), time.Second*3)
		stage        = actor.WithSystem(root)
	)

	defer cancel()
	stage.ActorOf(actor.FromProducer(NewTestActor))

	time.Sleep(time.Millisecond * 10)

	var i int32
	for i = 100; i < 110; i++ {
		test_client(stage, i)
	}

	stage.Wait()
	fmt.Println("结束")
}
