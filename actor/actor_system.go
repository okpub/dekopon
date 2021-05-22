package actor

import (
	"context"
	"sync"

	"github.com/okpub/dekopon/utils"
)

type actorSystem struct {
	Node
	utils.TaskDone
	ctx context.Context

	mu sync.Mutex
	PIDSet
}

func NewSystem() ActorSystem {
	return WithSystem(context.Background())
}

/*
* 自定义系统Actor
* 注意: 最好不要在退出actorSystem之前关闭ctx，而是调用Shutdown之后再关闭ctx
 */
func WithSystem(parent context.Context) ActorSystem {
	var (
		exit        = utils.MakeDone()
		ctx, cancel = context.WithCancel(parent)
		stage       = &actorSystem{TaskDone: exit, ctx: ctx, PIDSet: make(PIDSet)}
	)

	go func() {
		select {
		case <-ctx.Done():
			exit.Shutdown()
		case <-exit:
			cancel()
		}
	}()

	return stage.init()
}

//init stage
func (stage *actorSystem) init() ActorSystem {
	stage.system = stage
	return stage
}

//override
func (stage *actorSystem) ActorOf(props *Props, args ...ActorOption) (pid PID) {
	pid = props.spawn(stage, NewOptions(args))

	stage.mu.Lock()
	stage.PIDSet.Set(pid)
	stage.mu.Unlock()
	return
}

//system
func (stage *actorSystem) Wait() {
	utils.WaitDone(stage.Done())

	stage.mu.Lock()
	var list = stage.PIDSet.Values()
	stage.mu.Unlock()

	for _, pid := range list {
		pid.GracefulStop()
	}
}
