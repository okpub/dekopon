package actor

import (
	"context"
	"sync"
)

type WaitGroup struct {
	wg sync.WaitGroup
}

func (w *WaitGroup) Wrap(fn func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		fn()
	}()
}

func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

//static func
func Wait(ctx context.Context) {
	WaitDone(ctx.Done())
}

func WaitDone(exit <-chan struct{}) {
	<-exit
}
