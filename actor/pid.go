package actor

import (
	"context"

	"github.com/okpub/dekopon/utils"
)

type PID interface {
	//data
	Options() PIDOptions
	Background() context.Context

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
	*PIDOptions
	ctx context.Context
	ref ActorProcess
}

func NewPIDWithOptions(ctx context.Context, ref ActorProcess, opts *PIDOptions) PID {
	return &actorRef{ctx: ctx, ref: ref, PIDOptions: opts}
}

func NewPID(ctx context.Context, ref ActorProcess, args ...PIDOption) PID {
	return NewPIDWithOptions(ctx, ref, NewOptions(args))
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
	var options = NewPublish(message)
	options.Filler(args)
	return pid.ref.CallUserMessage(options.Context, pid, options.Message)
}

func (pid *actorRef) Send(message interface{}, args ...PublishOption) (err error) {
	var options = NewPublish(message)
	options.Filler(args)
	err = pid.sendUserMessage(options.Context, options.Message)
	return
}

func (pid *actorRef) Request(message interface{}, sender PID, args ...PublishOption) (err error) {
	var options = NewPublish(message)
	options.Filler(args)
	err = pid.sendUserMessage(options.Context, REQ(options.Message, sender))
	return
}

//stop
func (pid *actorRef) Close() (err error) {
	err = pid.ref.PostStop(pid)
	return
}

func (pid *actorRef) GracefulStop() (err error) {
	err = pid.ref.PostStop(pid)
	utils.Wait(pid.Background())
	return
}

//info
func (pid *actorRef) Background() context.Context {
	return pid.ctx
}
