package actor

import "context"

//同步消息（必定有一个回执）
type SyncMessageProcess struct {
	UntypeProcess
	Request
}

func (ref *SyncMessageProcess) PostUserMessage(ctx context.Context, pid PID, message interface{}) (err error) {
	err = ref.Respond(message)
	return
}

func (ref *SyncMessageProcess) PostStop(pid PID) (done <-chan struct{}, err error) {
	ref.Done()
	return
}

//同步pid
func NewSync(request Request) PID {
	return NewPID(&SyncMessageProcess{Request: request}, SetName("request"))
}
