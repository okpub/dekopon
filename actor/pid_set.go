package actor

import (
	"fmt"
	"sync"
)

type PIDSet map[PID]struct{}

func (hset PIDSet) Set(pid PID) PID {
	if hset == nil {
		fmt.Println("WARNING: the pidset is nil!")
	} else {
		hset[pid] = struct{}{}
	}
	return pid
}

func (hset PIDSet) Remove(pid PID) PID {
	delete(hset, pid)
	return pid
}

func (hset PIDSet) Each(fn func(PID)) {
	for pid := range hset {
		fn(pid)
	}
}

func (hset PIDSet) SafeEach(mu *sync.Mutex, fn func(PID)) {
	var list []PID
	mu.Lock()
	for pid := range hset {
		list = append(list, pid)
	}
	mu.Unlock()
	for _, pid := range list {
		fn(pid)
	}
}

func (hset PIDSet) Values() (arr []PID) {
	for pid := range hset {
		arr = append(arr, pid)
	}
	return
}
