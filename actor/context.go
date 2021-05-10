package actor

import (
	"context"
	"fmt"
	"time"
)

type actorContext struct {
	infoPart

	props *Props

	messageOrEnvelope interface{}
	actor             Actor

	extras     *actorContextExtras
	receiveDur time.Duration
}

//依赖父节点环境
func NewContext(info infoPart, props *Props) *actorContext {
	return &actorContext{infoPart: info, props: props}
}

//独立的context(不存在stage和parent)
func NewSelf(pid PID, props *Props) *actorContext {
	return &actorContext{infoPart: &Node{self: pid}, props: props}
}

func (ctx *actorContext) init() {
	if ctx.actor = ctx.props.NewActor(); ctx.actor == nil {
		ctx.actor = defaultEmptyActor
	}
}

func (ctx *actorContext) getExtras() *actorContextExtras {
	if ctx.extras == nil {
		var ttx ActorContext = ctx
		if chain := ctx.props.contextDecoratorChain; chain != nil {
			ttx = chain(ttx)
		}
		ctx.extras = newActorContextExtras(ttx)
	}
	return ctx.extras
}

//base
func (ctx *actorContext) Watch(who PID) error {
	return who.sendSystemMessage(&Watch{Who: ctx.Self()})
}

func (ctx *actorContext) Unwatch(who PID) error {
	return who.sendSystemMessage(&Unwatch{Who: ctx.Self()})
}

func (ctx *actorContext) Children() []PID {
	return ctx.getExtras().children.Values()
}

func (ctx *actorContext) ReceiveTimeout() time.Duration {
	return ctx.receiveDur
}

func (ctx *actorContext) SetReceiveTimeout(dur time.Duration) {
	ctx.receiveDur = dur
}

func (ctx *actorContext) CancelReceiveTimeout() {
	ctx.receiveDur = 0
}

//swapn
func (ctx *actorContext) ActorOf(props *Props, args ...PIDOption) PID {
	return ctx.getExtras().addChild(props.spawn(ctx, NewOptions(args)))
}

func (ctx *actorContext) Background() context.Context {
	return ctx.Self().Background()
}

//env
func (ctx *actorContext) Message() interface{} {
	return GetMessage(ctx.messageOrEnvelope)
}

func (ctx *actorContext) Sender() PID {
	return GetSender(ctx.messageOrEnvelope)
}

//invoker
func (ctx *actorContext) InvokeUserMessage(message interface{}) {
	switch event := message.(type) {
	case SystemMessage:
		ctx.InvokeSystemMessage(message)
	case func():
		ctx.handleFunc(event)
	default:
		ctx.processMessage(message)
	}
}

func (ctx *actorContext) InvokeSystemMessage(message interface{}) {
	switch event := message.(type) {
	case *Terminated:
		ctx.handleTerminated(event)
	case *RanError:
		ctx.handleFailure(event)
	case *Started:
		ctx.handleStart(event)
	case *Stopped:
		ctx.handleStop(event)
	case *Restarting:
		ctx.handleRestart(event)
	case *Watch:
		ctx.handleWatch(event)
	case *Unwatch:
		ctx.handleUnwatch(event)
	default:
		//fmt.Println("##CTX WARN: Miss handle [class SystemMessage]:", event)
		ctx.processMessage(message)
	}
}

func (ctx *actorContext) EscalateFailure(err error, message interface{}) {
	ctx.fireParent(&RanError{Err: err, Result: message, Who: ctx.Self()})
}

//private func
func (ctx *actorContext) fireParent(message interface{}) (err error) {
	if parent := ctx.Parent(); nil != parent {
		err = parent.sendSystemMessage(message)
	} else {
		err = fmt.Errorf("##GL fireParent: parent is nil!")
	}
	return
}
