package utils

type TaskDone chan struct{}

func MakeDone() TaskDone {
	return make(TaskDone)
}

func (done TaskDone) Done() <-chan struct{} {
	return done
}

func (done TaskDone) Shutdown() {
	SafeDone(done)
}
