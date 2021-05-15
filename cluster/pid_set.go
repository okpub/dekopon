package cluster

import "github.com/okpub/dekopon/actor"

type PIDSet map[string]actor.PID

func (childs PIDSet) Bind(name string, pid actor.PID) {
	childs[name] = pid
}

func (childs PIDSet) Get(name string) (pid actor.PID, ok bool) {
	pid, ok = childs[name]
	return
}
