package network

import (
	"fmt"
	"sync"

	"github.com/okpub/dekopon/actor"
)

type ServerProcess struct {
	actor.UntypeProcess
	server Server
	done   <-chan struct{}

	mu      sync.Mutex
	wathers actor.PIDSet
}

func NewServerProcess(done <-chan struct{}, server Server) actor.ActorProcess {
	return &ServerProcess{server: server, done: done, wathers: make(actor.PIDSet)}
}

func (p *ServerProcess) PostSystemMessage(pid actor.PID, message interface{}) (err error) {
	switch event := message.(type) {
	case *actor.Watch:
		p.handleWatch(event)
	case *actor.Unwatch:
		p.handleUnwatch(event)
	case *actor.Started:
		p.handleStart(pid)
	case *actor.Stopped:
		p.handleStop(pid)
	default:
		fmt.Println("Miss system-message:", message)
	}
	return
}

func (p *ServerProcess) PostStop(pid actor.PID) (done <-chan struct{}, err error) {
	done, err = p.done, p.server.Close()
	return
}

//handler
func (p *ServerProcess) handleWatch(event *actor.Watch) {
	p.mu.Lock()
	p.wathers.Set(event.Who)
	p.mu.Unlock()
}

func (p *ServerProcess) handleUnwatch(event *actor.Unwatch) {
	p.mu.Lock()
	p.wathers.Remove(event.Who)
	p.mu.Unlock()
}

func (p *ServerProcess) handleStart(pid actor.PID) {
	//p.server.ListenAndServe(pid.Options().Context, p.handler)
}

func (p *ServerProcess) handleStop(self actor.PID) {
	p.mu.Lock()
	var wathcers = p.wathers.Values()
	p.mu.Unlock()

	for _, pid := range wathcers {
		pid.Send(&actor.Terminated{Who: self})
	}
}
