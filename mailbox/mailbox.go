package mailbox

import (
	"context"

	"github.com/okpub/dekopon/utils"
)

type Producer func() Mailbox

type Mailbox interface {
	Start()
	RegisterHander(InvokerMessage, Dispatcher)

	CallUserMessage(context.Context, interface{}) (interface{}, error)
	PostUserMessage(context.Context, interface{}) error
	PostSystemMessage(interface{}) error

	Close() error
}

type InvokerMessage interface {
	InvokeUserMessage(interface{})
	InvokeSystemMessage(interface{})
	EscalateFailure(error, interface{})
}

//默认邮箱
type defaultMailbox struct {
	taskMailbox utils.TaskBuffer

	invoker    InvokerMessage
	dispatcher Dispatcher
}

func (box *defaultMailbox) Start() {
	var (
		readChan = box.taskMailbox.Recv()
	)
	for message := range readChan {
		box.invoker.InvokeUserMessage(message)
	}
}

func (box *defaultMailbox) RegisterHander(invoker InvokerMessage, dispatcher Dispatcher) {
	box.invoker = invoker
	box.dispatcher = dispatcher
}

func (box *defaultMailbox) PostSystemMessage(message interface{}) error {
	return box.taskMailbox.Send(message)
}

func (box *defaultMailbox) PostUserMessage(ctx context.Context, message interface{}) (err error) {
	err = box.taskMailbox.SendCtx(ctx, message)
	return
}

func (box *defaultMailbox) CallUserMessage(ctx context.Context, message interface{}) (interface{}, error) {
	var request = NewRequest(message)
	if err := box.PostUserMessage(ctx, request); utils.Die(err) {
		request.Respond(err)
	}
	return request.Body(ctx)
}

func (box *defaultMailbox) Close() error {
	return box.taskMailbox.Close()
}
