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
type Test struct {
}

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
	case actor.Request:
		if req, ok := event.Message().(*rpc.Request); ok {
			a.processRequest(ctx, event, req)
		}
	}
}

func (a *TestCallActor) processRequest(ctx actor.ActorContext, event actor.Request, req *rpc.Request) {
	var res, err = ValueOf(a.Manager.Router(req.ServerName, req.MethodName, req))
	if err == nil {
		event.Respond(res)
	} else {
		event.Respond(err)
	}
}

func rpc_client() {
	var client = NewClient(network.SetDialAddr("localhost:9098"))
	client.Start(context.Background())

	for i := 0; i < 10; i++ {
		var t = time.Now()
		var res, _ = client.Request(context.Background(), &rpc.Request{ServerName: "Test", MethodName: "Login"})
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
