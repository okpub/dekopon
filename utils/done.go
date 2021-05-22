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

func SafeDone(done chan<- struct{}) (err error) {
	defer func() { err = CatchDie(recover()) }()
	close(done)
	return
}
