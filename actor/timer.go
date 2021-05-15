package actor

import (
	"context"
	"time"
)

func OnceTimer(pid PID, d time.Duration, fn func()) (cancel context.CancelFunc) {
	var (
		child context.Context
		timer = time.NewTimer(d)
	)

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

func LoopTimer(pid PID, d time.Duration, fn func()) (cancel context.CancelFunc) {
	var (
		child context.Context
		timer = time.NewTicker(d)
	)

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
