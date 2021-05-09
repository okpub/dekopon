package actor

func (ctx *actorContext) handleFunc(event func()) {
	event()
}

func (ctx *actorContext) handleStart(event *Started) {
	ctx.init()
	ctx.processMessage(event)
}

func (ctx *actorContext) handleStop(event *Stopped) {
	//stop and wait childs
	ctx.getExtras().stopAllChildren()
	//handler message
	ctx.processMessage(event)
	//talk watchers
	ctx.getExtras().stopBroadcast(ctx.Self())
	//send to parent
	ctx.fireParent(&Terminated{Who: ctx.Self()})
}

func (ctx *actorContext) handleRestart(event *Restarting) {
	//1 stop childs
	ctx.getExtras().stopAllChildren()
	//2 init ctx
	ctx.init()
	//3 handler event
	ctx.processMessage(event)
}

func (ctx *actorContext) handleWatch(event *Watch) {
	ctx.getExtras().watch(event.GetWho())
}

func (ctx *actorContext) handleUnwatch(event *Unwatch) {
	ctx.getExtras().unwatch(event.GetWho())
}

func (ctx *actorContext) handleFailure(event *RanError) {
	ctx.processMessage(event)
}

func (ctx *actorContext) handleTerminated(event *Terminated) {
	//1 remove child
	ctx.getExtras().removeChild(event.GetWho())
	//2 handler event
	ctx.processMessage(event)
}
