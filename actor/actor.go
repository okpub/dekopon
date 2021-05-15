package actor

import (
	"time"
)

type Actor interface {
	Received(ActorContext)
}

type Producer func() Actor

//静态函数actor
type ActorFunc func(ActorContext)

var defaultEmptyActor = ActorFunc(func(ctx ActorContext) {})

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
		messagePart
	}

	ReceiverContext interface {
		infoPart
		recvPart
		messagePart
	}

	SenderContext interface {
		infoPart
		sendPart
		messagePart
	}
)

//part
type (
	infoPart interface {
		ChildOf(PID) infoPart
		Self() PID
		Parent() PID
		System() ActorSystem
	}

	messagePart interface {
		Message() interface{}
	}

	spawnPart interface {
		ActorOf(*Props, ...ActorOption) PID
	}

	recvPart interface {
		Received(MessageEnvelope)
	}

	sendPart interface {
		Sender() PID
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
		messagePart
		Sender() PID
	}

	Request interface {
		messagePart
		Done()
		Respond(interface{}) error
	}
)
