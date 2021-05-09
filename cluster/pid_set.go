package cluster

import "github.com/skimmer/actor"

type PIDSet map[string]actor.PID

func (childs PIDSet) Bind(name string, pid actor.PID) {
	childs[name] = pid
}
