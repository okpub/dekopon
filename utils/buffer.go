package utils

import (
	"context"
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

func (ch TaskBuffer) Recv() <-chan interface{} {
	return ch
}

func (ch TaskBuffer) Send(message interface{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	ch <- message
	return
}

//当前堵塞立即返回
func (ch TaskBuffer) SendTry(message interface{}) (err error) {
	defer func() {
		if err == nil {
			err = CatchDie(recover())
		}
	}()

	select {
	case ch <- message:
		//send ok
	default:
		err = TempErr
	}
	return
}

func (ch TaskBuffer) SendCtx(ctx context.Context, message interface{}) (err error) {
	defer func() {
		if err == nil {
			err = CatchDie(recover())
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

func (ch TaskBuffer) Close() (err error) {
	err = SafeClose(ch)
	return
}
