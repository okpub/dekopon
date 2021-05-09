package mailbox

import (
	"context"
	"fmt"
	"time"

	"github.com/skimmer/utils"
)

const (
	Blocking          = 0
	UnBlocking        = -1
	DefaultPendingNum = 20
)

type TaskBuffer chan interface{}

func MakeBuffer(size int) TaskBuffer {
	return make(TaskBuffer, size)
}

func MakeBlock() TaskBuffer {
	return make(TaskBuffer)
}

func MakeTest() TaskBuffer {
	return make(TaskBuffer, DefaultPendingNum)
}

func (ch TaskBuffer) Start(_ context.Context) error {
	panic(fmt.Errorf("untype func Start!"))
}

func (ch TaskBuffer) RegisterHander(_ InvokerMessage) {
	panic(fmt.Errorf("untype func RegisterHander!"))
}

func (ch TaskBuffer) Recv() <-chan interface{} {
	return ch
}

func (ch TaskBuffer) Send(message interface{}) (err error) {
	defer func() { err = utils.CatchDie(recover()) }()
	ch <- message
	return
}

func (ch TaskBuffer) SendTimeout(message interface{}, dur time.Duration) (err error) {
	defer func() {
		if err == nil {
			err = utils.CatchDie(recover())
		}
	}()

	switch {
	case dur == Blocking:
		ch <- message
	case dur < Blocking:
		select {
		case ch <- message:
			//send ok
		default:
			err = TempErr
		}
	default:
		select {
		case ch <- message:
			//send ok
		case <-time.After(dur):
			err = TempErr
		}
	}
	return
}

func (ch TaskBuffer) CallUserMessage(ctx context.Context, message interface{}) (interface{}, error) {
	var request = NewRequest(message)
	if err := ch.PostUserMessage(ctx, request); Fail(err) {
		request.Respond(err)
	}
	return request.Body(ctx)
}

func (ch TaskBuffer) PostUserMessage(ctx context.Context, message interface{}) (err error) {
	defer func() {
		if err == nil {
			err = utils.CatchDie(recover())
		}
	}()
	select {
	case ch <- message:
		//send ok
	case <-ctx.Done():
		err = ctx.Err()
	}
	return
}

func (ch TaskBuffer) PostSystemMessage(message interface{}) error {
	return ch.Send(message)
}

//read
func (ch TaskBuffer) Read(ctx context.Context) (message interface{}, err error) {
	var ok = false
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case message, ok = <-ch:
		if ok {
			//read ok
		} else {
			err = utils.EOF
		}
	}
	return
}

func (ch TaskBuffer) Close() (err error) {
	err = utils.SafeClose(ch)
	return
}
