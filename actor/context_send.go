package actor

import (
	"context"
	"fmt"
)

/*
* ctx发送的用户消息不会是不会超时的
 */
func (ctx *actorContext) sendUserMessage(pid PID, message interface{}) (err error) {
	if pid == nil {
		err = fmt.Errorf("##CTX sendUserMessage: pid is nil")
	} else {
		var mctx = ctx.getExtras().context
		if chain := ctx.props.senderMiddlewareChain; nil != chain {
			err = chain(context.Background(), mctx, pid, WrapEnvelope(message))
		} else {
			err = pid.sendUserMessage(context.Background(), message)
		}
	}
	return
}

//send part
func (ctx *actorContext) Send(pid PID, message interface{}) (err error) {
	return ctx.sendUserMessage(pid, message)
}

func (ctx *actorContext) RequestWithCustomSender(pid PID, message interface{}, sender PID) (err error) {
	return ctx.sendUserMessage(pid, REQ(message, sender))
}

//send other
func (ctx *actorContext) Respond(message interface{}) error {
	var mctx = ctx.getExtras().context
	return mctx.Send(mctx.Sender(), message)
}

func (ctx *actorContext) Bubble(message interface{}) error {
	var mctx = ctx.getExtras().context
	return mctx.Send(mctx.Parent(), message)
}

func (ctx *actorContext) Forward(pid PID) (err error) {
	var mctx = ctx.getExtras().context
	switch ctx.messageOrEnvelope.(type) {
	case SystemMessage:
		err = fmt.Errorf("##CTX Forward Err: can't forward type=%v", ctx.messageOrEnvelope)
	default:
		err = mctx.Send(pid, ctx.messageOrEnvelope)
	}
	return
}

func (ctx *actorContext) Request(pid PID, message interface{}) error {
	var mctx = ctx.getExtras().context
	return mctx.RequestWithCustomSender(pid, message, mctx.Self())
}

//stop
func (ctx *actorContext) Stop(pid PID) (err error) {
	if nil == pid {
		err = fmt.Errorf("##CTX Stop Err: pid is nil")
	} else {
		err = pid.Close()
	}
	return
}
