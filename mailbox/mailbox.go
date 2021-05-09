package mailbox

import (
	"context"
)

func Fail(err error) bool {
	return err != nil
}

type Producer func() Mailbox

type Mailbox interface {
	Start(context.Context) error
	RegisterHander(InvokerMessage)
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

type MessageRespond interface {
	Body(context.Context) (interface{}, error)
}

type MessageRequest interface {
	Message() interface{}      //请求消息
	Respond(interface{}) error //主动回答
	Done()
}
