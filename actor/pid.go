package actor

import (
	"context"

	"github.com/okpub/dekopon/utils"
)

type PID interface {
	//data
	Options() ActorOptions

	//private send
	sendSystemMessage(interface{}) error
	sendUserMessage(context.Context, interface{}) error

	//send
	Call(interface{}, ...PublishOption) (interface{}, error)
	Send(interface{}, ...PublishOption) error
	Request(interface{}, PID, ...PublishOption) error

	//stopped
	Close() error
	GracefulStop() error
}

//class pid
type actorRef struct {
	*ActorOptions
	ref ActorProcess
}

func WithPID(ref ActorProcess, opts *ActorOptions) PID {
	return &actorRef{ref: ref, ActorOptions: opts}
}

func NewPID(ref ActorProcess, args ...ActorOption) PID {
	return WithPID(ref, NewOptions(args))
}

//private
func (pid *actorRef) sendSystemMessage(message interface{}) error {
	return pid.ref.PostSystemMessage(pid, message)
}

func (pid *actorRef) sendUserMessage(ctx context.Context, message interface{}) error {
	return pid.ref.PostUserMessage(ctx, pid, message)
}

//public
func (pid *actorRef) Call(message interface{}, args ...PublishOption) (interface{}, error) {
	var options = NewPublish(message, args...)
	return pid.ref.CallUserMessage(options.Context, pid, options.Message)
}

func (pid *actorRef) Send(message interface{}, args ...PublishOption) (err error) {
	var options = NewPublish(message, args...)
	err = pid.sendUserMessage(options.Context, options.Message)
	return
}

func (pid *actorRef) Request(message interface{}, sender PID, args ...PublishOption) (err error) {
	var options = NewPublish(message, args...)
	err = pid.sendUserMessage(options.Context, WrapMessage(options.Message, sender))
	return
}

//stop
func (pid *actorRef) Close() (err error) {
	_, err = pid.ref.PostStop(pid)
	return
}

func (pid *actorRef) GracefulStop() error {
	var done, err = pid.ref.PostStop(pid)
	utils.WaitDone(done)
	return err
}
