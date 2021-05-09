package actor

import (
	"context"

	"github.com/okpub/dekopon/mailbox"
)

/*
* 默认的消息邮箱
 */
type actorMailbox struct {
	mailbox.TaskBuffer
	invoker mailbox.InvokerMessage
}

func NewMailbox() mailbox.Mailbox {
	return &actorMailbox{
		TaskBuffer: mailbox.MakeBuffer(defaultPendingNum),
	}
}

func (this *actorMailbox) RegisterHander(invoker mailbox.InvokerMessage) {
	this.invoker = invoker
}

func (this *actorMailbox) Start(ctx context.Context) (err error) {
	this.invoker.InvokeSystemMessage(EVENT_START)
	defer this.invoker.InvokeSystemMessage(EVENT_STOP)
	for message := range this.TaskBuffer {
		this.invoker.InvokeUserMessage(message)
	}
	return
}
