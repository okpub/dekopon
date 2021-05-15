package cluster

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/network"
	"github.com/okpub/dekopon/observer"
)

//最终处理
type Test struct{}

func (t *Test) Login(req *rpc.Request) (res *rpc.Response, err error) {
	fmt.Println("登陆")
	return
}

//同步处理模型(转接口)
type TestCallActor struct {
	*observer.Manager
}

func NewTestCallActor() actor.Actor {
	var ob = observer.NewManager()
	ob.Register(&Test{})
	return &TestCallActor{
		Manager: ob,
	}
}

func (a *TestCallActor) Received(ctx actor.ActorContext) {
	switch event := ctx.Message().(type) {
	case *actor.Started:
		//todo
	case *rpc.Request:
		a.processRequest(ctx, event)

	}
}

func (a *TestCallActor) processRequest(ctx actor.ActorContext, event *rpc.Request) {
	var res, err = ValueOf(a.Manager.Router(event.ServerName, event.MethodName, event))
	if err == nil {
		ctx.Respond(res)
	} else {
		ctx.Respond(err)
	}
}

func rpc_client() {
	var client = NewClient(network.SetDialAddr("localhost:9098"))
	client.Start()

	for i := 0; i < 10; i++ {
		var t = time.Now()
		var res, _ = client.Request(&rpc.Request{ServerName: "Test", MethodName: "Login"})
		fmt.Println("得到数据:", res, time.Since(t))
	}
}

func TestRpcInit(t *testing.T) {
	var root, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var stage = actor.WithSystem(root)

	var pid = stage.ActorOf(actor.FromProducer(NewTestCallActor))
	var p = make(PIDSet)
	p.Bind("Test", pid)

	var svr = NewRPCServer(p, network.SetAddr(":9098"))
	go svr.Serve(root)

	time.Sleep(time.Millisecond * 10)
	rpc_client()

	stage.Wait()
	fmt.Println("关闭")
}
