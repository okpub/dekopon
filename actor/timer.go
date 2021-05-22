package actor

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	time_count  int64
	time_active int64
)

func OnceTimer(pid PID, dur time.Duration, fn func()) (cancel context.CancelFunc) {
	var (
		child context.Context
		timer = time.NewTimer(dur)
		id    = atomic.AddInt64(&time_count, 1)
	)

	atomic.AddInt64(&time_active, 1)
	child, cancel = context.WithCancel(context.Background())

	var method = func() {
		select {
		case <-child.Done():
			//timer closed
		default:
			cancel()
			fn()
		}
	}

	go func() {
		defer atomic.AddInt64(&time_active, -1)
		defer fmt.Println("Error: time stop once id=", id)
		select {
		case <-timer.C:
			if err := pid.Send(method); err == nil {
				//send ok
			} else {
				cancel()
			}
		case <-child.Done():
			timer.Stop()
		}
	}()

	return
}

func LoopTimer(pid PID, dur time.Duration, fn func()) (cancel context.CancelFunc) {
	var (
		child context.Context
		timer = time.NewTicker(dur)
		id    = atomic.AddInt64(&time_count, 1)
	)

	atomic.AddInt64(&time_active, 1)
	child, cancel = context.WithCancel(context.Background())
	var method = func() {
		select {
		case <-child.Done():
			//timer closed
		default:
			fn()
		}
	}

	go func() {
		defer atomic.AddInt64(&time_active, -1)
		defer fmt.Println("Error: time stop loop id=", id)
		for {
			select {
			case <-timer.C:
				if err := pid.Send(method); err == nil {
					//send ok
				} else {
					cancel()
				}
			case <-child.Done():
				timer.Stop()
				return
			}
		}
	}()

	return
}

func TimerCount() int64 {
	return atomic.LoadInt64(&time_active)
}
