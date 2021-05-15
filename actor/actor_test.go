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
		fmt.Println("关闭:", ctx.Self())
	case int:
		ctx.Respond("傻逼")
	case string:
		time.Sleep(time.Millisecond * 100)
		if str == "shutdown" {
			ctx.System().Shutdown()
		} else if str == "child" {
			var pid = ctx.ActorOf(From(handlerMessage))
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
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
		stage       = WithSystem(context.WithValue(ctx, "router", "路由器"))
		pid         = stage.ActorOf(From(handlerMessage))
	)
	defer cancel()

	pid.Send("child")
	var res, err = pid.Call(13)
	fmt.Println(res, err)
	pid.Send("shutdown")

	stage.Wait()
	fmt.Println("最终推出")
}
