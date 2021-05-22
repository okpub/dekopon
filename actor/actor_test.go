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
		OnceTimer(ctx.Self(), time.Millisecond*10, func() {
			fmt.Println("计时器内部执行")
			fmt.Println("计时器个数:", TimerCount())
		})
		fmt.Println("计时器个数2:", TimerCount())
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
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		stage       = WithSystem(ctx)
		pid         = stage.ActorOf(From(handlerMessage))
	)
	defer cancel()

	pid.Send("child1")
	pid.Send("child2")
	pid.Send("child3")

	var res, err = pid.Call(13)
	fmt.Println(res, err)
	pid.Send("shutdown")

	stage.Wait()
	fmt.Println("最终推出")
}
