package mailbox

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func handlerMessage(message interface{}) {
	if request, ok := message.(*mailboxRequest); ok {
		defer request.Done()
		request.Respond("没有消费")
	} else {
		fmt.Println("消费:", message)
	}
}

func TestInit(t *testing.T) {
	var (
		buff = NewMailbox()
	)

	buff.RegisterHander(InvokerFunc(handlerMessage), NewDefaultDispatcher())
	go buff.Start()

	buff.PostSystemMessage("start")
	buff.PostUserMessage(context.Background(), "为啥")
	var resp, err = buff.CallUserMessage(context.Background(), "我是会")
	fmt.Println(resp, err)

	time.Sleep(time.Second)
	buff.Close()
}
