package mailbox

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func handlerMessage(message interface{}) {
	if request, ok := message.(MessageRequest); ok {
		defer request.Done()
		request.Respond("没有消费")
	} else {
		fmt.Println("消费:", message)
	}
}

func TestInit(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		buff        = MakeBuffer(10)
	)

	go func() {
		defer cancel()
		for message := range buff {
			handlerMessage(message)
		}
	}()

	buff.Send("为啥")
	var resp, err = buff.CallUserMessage(ctx, "我是会")
	fmt.Println(resp, err)

	time.Sleep(time.Second)
	buff.Close()
	<-ctx.Done()
}
