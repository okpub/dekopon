package actor

import (
	"context"
	"time"
)

type Actor interface {
	Received(ActorContext)
}

type Producer func() Actor

//静态函数actor
type ActorFunc func(ActorContext)

func (method ActorFunc) Received(ctx ActorContext) {
	method(ctx)
}

//context
type ActorSystem interface {
	spawnPart
	infoPart
	Shutdown()
	Wait()
}

type (
	SpawnContext interface {
		infoPart
		spawnPart
	}

	ActorContext interface {
		basePart
		infoPart
		spawnPart
		recvPart
		sendPart
		stopPart
	}

	ReceiverContext interface {
		sendPart
		recvPart
	}

	SenderContext interface {
		sendPart
	}
)

//part
type (
	infoPart interface {
		Self() PID
		Parent() PID
		System() ActorSystem
		ChildOf(PID) infoPart
	}

	spawnPart interface {
		ActorOf(*Props, ...PIDOption) PID
		Background() context.Context
	}

	recvPart interface {
		MessageEnvelope
		Received(MessageEnvelope)
	}

	sendPart interface {
		Send(PID, interface{}) error
		RequestWithCustomSender(PID, interface{}, PID) error
		//other
		Bubble(interface{}) error
		Request(PID, interface{}) error
		Respond(interface{}) error
		Forward(PID) error
	}

	stopPart interface {
		Stop(PID) error
	}

	basePart interface {
		Watch(PID) error
		Unwatch(PID) error
		Children() []PID
		SetReceiveTimeout(time.Duration)
		CancelReceiveTimeout()
		ReceiveTimeout() time.Duration
		//		Respond(interface{}) error
		//		Forward(PID) error
	}
)

//message
type (
	MessageEnvelope interface {
		Sender() PID
		Message() interface{}
	}

	Request interface {
		Done()
		Message() interface{}
		Respond(interface{}) error
	}
)
