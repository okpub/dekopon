package actor

/*
* 真正的核心
 */
type actorContextExtras struct {
	children PIDSet
	watchers PIDSet
	context  ActorContext
}

func newActorContextExtras(context ActorContext) *actorContextExtras {
	return &actorContextExtras{
		context:  context,
		children: make(PIDSet),
		watchers: make(PIDSet),
	}
}

func (ctxExt *actorContextExtras) addChild(pid PID) PID {
	ctxExt.children.Set(pid)
	return pid
}

func (ctxExt *actorContextExtras) removeChild(pid PID) {
	ctxExt.children.Remove(pid)
}

func (ctxExt *actorContextExtras) watch(watcher PID) {
	ctxExt.watchers.Set(watcher)
}

func (ctxExt *actorContextExtras) unwatch(watcher PID) {
	ctxExt.watchers.Remove(watcher)
}

func (ctxExt *actorContextExtras) stopAllChildren() {
	var children = ctxExt.children
	ctxExt.children = make(PIDSet)
	children.Each(func(pid PID) { pid.GracefulStop() })
}

func (ctxExt *actorContextExtras) stopBroadcast(self PID) {
	ctxExt.watchers.Each(func(pid PID) {
		pid.sendSystemMessage(&Terminated{Who: self})
	})
}
