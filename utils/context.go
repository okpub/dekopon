package utils

import (
	"context"
)

//static func
func Wait(ctx context.Context) {
	WaitDone(ctx.Done())
}

func WaitDone(exit <-chan struct{}) {
	<-exit
}
