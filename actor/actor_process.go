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
	PostStop(PID) (<-chan struct{}, error)
}

//默认
type defaultProcess struct {
	done    <-chan struct{}
	mailbox mailbox.Mailbox
}

func NewDfaultProcess(done <-chan struct{}, mailbox mailbox.Mailbox) ActorProcess {
	return &defaultProcess{done: done, mailbox: mailbox}
}

func (ref *defaultProcess) CallUserMessage(ctx context.Context, pid PID, message interface{}) (interface{}, error) {
	return ref.mailbox.CallUserMessage(ctx, message)
}

func (ref *defaultProcess) PostUserMessage(ctx context.Context, pid PID, message interface{}) error {
	return ref.mailbox.PostUserMessage(ctx, message)
}

func (ref *defaultProcess) PostSystemMessage(pid PID, message interface{}) error {
	return ref.mailbox.PostSystemMessage(message)
}

func (ref *defaultProcess) PostStop(pid PID) (done <-chan struct{}, err error) {
	done, err = ref.done, ref.mailbox.Close()
	return
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

func (*UntypeProcess) PostStop(pid PID) (<-chan struct{}, error) {
	panic(fmt.Errorf("[Class UntypeProcess] unrealized PostStop"))
}
