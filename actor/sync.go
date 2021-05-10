package actor

import (
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
