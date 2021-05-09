package actor

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func handlerMessage(ctx ActorContext) {
	switch str := ctx.Message().(type) {
	case *Stopped:
		fmt.Println("关闭:", ctx.Background().Value("router"))
	case string:
		time.Sleep(time.Millisecond * 100)
		if str == "shutdown" {
			ctx.System().Shutdown()
		} else if str == "child" {
			var pid = ctx.ActorOf(From(handlerMessage).WithValue("router", "代替陆游"))
			pid.Send("发送")
			pid.Send("my self")
			var err = pid.Send("我来发消息", SetTimeout(time.Millisecond*10))
			fmt.Println(err)
		} else {
			fmt.Println("todo:", str)
		}
	}
}

func TestInit(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		stage       = WithSystem(context.WithValue(ctx, "router", "路由器"))
		pid         = stage.ActorOf(From(handlerMessage))
	)
	defer cancel()

	pid.Send("child")

	for i := 0; i < 10; i++ {
		//fmt.Println(pid.Call("sync message"))
		pid.Send("child")
	}

	pid.Send("shutdown")

	stage.Wait()
	time.Sleep(time.Second)
}
