package actor

func (ctx *actorContext) Received(env MessageEnvelope) {
	switch event := env.Message().(type) {
	case Request:
		ctx.defaultRequestReceive(event)
		ctx.finishReceive()
	default:
		ctx.defaultReceive(env)
		ctx.finishReceive()
	}
}

//1 异步消息
func (ctx *actorContext) defaultReceive(env MessageEnvelope) {
	ctx.messageOrEnvelope = env
	ctx.actor.Received(ctx.getExtras().context)
}

//2 同步请求
func (ctx *actorContext) defaultRequestReceive(request Request) {
	var sender = NewSync(request)
	ctx.defaultReceive(WrapMessage(request.Message(), sender))
	sender.Close()
}

func (ctx *actorContext) finishReceive() {
	ctx.messageOrEnvelope = nil
}
