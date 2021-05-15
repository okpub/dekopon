package actor

func (ctx *actorContext) Received(env MessageEnvelope) {
	if request, ok := env.Message().(Request); ok {
		ctx.defaultRequestReceive(request)
	} else {
		ctx.messageOrEnvelope = env
		ctx.defaultReceive()
		ctx.messageOrEnvelope = nil
	}
}

//private
func (ctx *actorContext) processMessage(message interface{}) {
	if chain := ctx.props.receiverMiddlewareChain; chain != nil {
		chain(ctx.getExtras().context, WrapEnvelope(message))
	} else {
		ctx.getExtras().context.Received(WrapEnvelope(message))
	}
}

func (ctx *actorContext) defaultReceive() {
	ctx.actor.Received(ctx.getExtras().context)
}

//处理同步消息
func (ctx *actorContext) defaultRequestReceive(request Request) {
	var sender = SyncPID(request)
	ctx.messageOrEnvelope = WrapMessage(request.Message(), sender)
	ctx.defaultReceive()
	sender.Close()
	ctx.messageOrEnvelope = nil
}
