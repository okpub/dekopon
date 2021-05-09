package actor

import (
	"context"

	"github.com/skimmer/utils"
)

type actorSystem struct {
	Node
	TaskDone
	WaitGroup
	ctx context.Context
}

func NewSystem() ActorSystem {
	return WithSystem(context.Background())
}

/*
* 自定义系统Actor
* 注意: 最好不要在退出actorSystem之前关闭ctx，而是调用Shutdown之后再关闭ctx
 */
func WithSystem(ctx context.Context) ActorSystem {
	var (
		exit  = MakeDone()
		stage = &actorSystem{TaskDone: exit, ctx: ctx}
	)

	go func() {
		select {
		case <-ctx.Done():
			utils.SafeDone(exit)
		case <-exit:
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
func (stage *actorSystem) ActorOf(props *Props, args ...PIDOption) (pid PID) {
	pid = props.spawn(stage, NewOptions(args))
	stage.Wrap(func() {
		select {
		case <-pid.Background().Done():
			//pid closed
		case <-stage.Done():
			pid.GracefulStop()
		}
	})

	return
}

func (stage *actorSystem) Background() context.Context {
	return stage.ctx
}

//system
func (stage *actorSystem) Shutdown() {
	stage.Close()
}

func (stage *actorSystem) Wait() {
	WaitDone(stage.Done())
	stage.WaitGroup.Wait()
}
