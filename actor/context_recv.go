package actor

func (ctx *actorContext) Received(env MessageEnvelope) {
	ctx.messageOrEnvelope = env
	ctx.defaultReceive()
	ctx.handleRequest(env.Message())
	ctx.messageOrEnvelope = nil
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

func (*actorContext) handleRequest(message interface{}) {
	if req, ok := message.(Request); ok {
		req.Done()
	}
}
