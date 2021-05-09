package actor

import (
	"context"
	"fmt"

	"github.com/okpub/dekopon/mailbox"
)

/*
* 处理pid过程
 */
type ActorProcess interface {
	CallUserMessage(context.Context, PID, interface{}) (interface{}, error)
	PostUserMessage(context.Context, PID, interface{}) error
	PostSystemMessage(PID, interface{}) error
	PostStop(PID) error
}

type actorProcess struct {
	mailbox mailbox.Mailbox
}

func NewProcess(mailbox mailbox.Mailbox) ActorProcess {
	return &actorProcess{mailbox: mailbox}
}

func (ref *actorProcess) CallUserMessage(ctx context.Context, pid PID, message interface{}) (interface{}, error) {
	return ref.mailbox.CallUserMessage(ctx, message)
}

func (ref *actorProcess) PostUserMessage(ctx context.Context, pid PID, message interface{}) error {
	return ref.mailbox.PostUserMessage(ctx, message)
}

func (ref *actorProcess) PostSystemMessage(pid PID, message interface{}) error {
	return ref.mailbox.PostSystemMessage(message)
}

func (ref *actorProcess) PostStop(pid PID) error {
	return ref.mailbox.Close()
}

//未实现过程
type UntypeProcess struct{}

func (*UntypeProcess) CallUserMessage(ctx context.Context, pid PID, message interface{}) (interface{}, error) {
	panic(fmt.Errorf("[Class UntypeProcess] unrealized CallUserMessage"))
}

func (*UntypeProcess) PostUserMessage(ctx context.Context, pid PID, message interface{}) error {
	panic(fmt.Errorf("[Class UntypeProcess] unrealized PostUserMessage"))
}

func (*UntypeProcess) PostSystemMessage(pid PID, message interface{}) error {
	panic(fmt.Errorf("[Class UntypeProcess] unrealized PostSystemMessage"))
}

func (*UntypeProcess) PostStop(pid PID) error {
	panic(fmt.Errorf("[Class UntypeProcess] unrealized PostStop"))
}
