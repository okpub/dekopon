package utils

import (
	"context"
	"time"
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

func (ch TaskBuffer) Close() error {
	return SafeClose(ch)
}

//funcs
func Send(ch chan<- interface{}, message interface{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	ch <- message
	return
}

//当前堵塞立即返回
func SendTry(ch chan<- interface{}, message interface{}) (err error) {
	defer func() {
		if err == nil {
			err = CatchDie(recover())
		}
	}()

	select {
	case ch <- message:
		//send ok
	default:
		err = TimeoutErr
	}
	return
}

func SendCtx(ch chan<- interface{}, ctx context.Context, message interface{}) (err error) {
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

func SendTimeout(ch chan<- interface{}, message interface{}, dur time.Duration) (err error) {
	defer func() {
		if err == nil {
			err = CatchDie(recover())
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
			err = TimeoutErr
		}
	default:
		select {
		case ch <- message:
			//send ok
		case <-time.After(dur):
			err = TimeoutErr
		}
	}
	return
}

func SafeClose(task chan interface{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	close(task)
	return
}
