package actor

import "github.com/skimmer/utils"

type TaskDone chan struct{}

func MakeDone() TaskDone {
	return make(TaskDone)
}

func (done TaskDone) Done() <-chan struct{} {
	return done
}

func (done TaskDone) Close() (err error) {
	err = utils.SafeDone(done)
	return
}
