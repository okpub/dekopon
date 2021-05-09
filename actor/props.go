package actor

import (
	"context"

	"github.com/skimmer/mailbox"
)

const defaultPendingNum = 20

var (
	defaultMailboxProducer = NewMailbox
	defaultSpawner         = func(parent SpawnContext, props *Props, options *PIDOptions) PID {
		var (
			child, cancel = props.WithCancel(parent.Background())
			box           = props.NewMailbox()
			pid           = NewPIDWithOptions(child, NewProcess(box), options)
			ctx           = NewContext(parent.ChildOf(pid), props)
		)

		box.RegisterHander(ctx)

		go func() {
			defer cancel()
			box.Start(child)
		}()
		return pid
	}
)

//目前够用
type Props struct {
	producer        Producer
	spawner         SpawnFunc
	mailboxProducer mailbox.Producer

	//value
	valueMiddleware      []ValueMiddleware
	valueMiddlewareChain ValueFunc

	//middkeware
	receiverMiddleware      []ReceiverMiddleware
	receiverMiddlewareChain ReceiverFunc
	senderMiddleware        []SenderMiddleware
	senderMiddlewareChain   SenderFunc
	spawnMiddleware         []SpawnMiddleware
	spawnMiddlewareChain    SpawnFunc
	contextDecorator        []ContextDecorator
	contextDecoratorChain   ContextDecoratorFunc
}

func (props *Props) NewActor() Actor {
	return props.producer()
}

func (props *Props) NewMailbox() mailbox.Mailbox {
	return props.mailboxProducer()
}

//context apply
func (props *Props) spawn(parent SpawnContext, options *PIDOptions) (pid PID) {
	if chain := props.spawnMiddlewareChain; chain != nil {
		return chain(parent, props, options)
	}
	return props.spawner(parent, props, options)
}

func (props *Props) WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	var ctx = parent
	if chain := props.valueMiddlewareChain; chain != nil {
		ctx = chain(parent)
	}
	return context.WithCancel(ctx)
}

//api with
func (props *Props) WithMailboxProducer(p mailbox.Producer) *Props {
	props.mailboxProducer = p
	return props
}

func (props *Props) WithProducer(p Producer) *Props {
	props.producer = p
	return props
}

func (props *Props) WithFunc(p ActorFunc) *Props {
	return props.WithProducer(func() Actor { return p })
}

//new props
func From(p ActorFunc) *Props {
	return FromProducer(func() Actor { return p })
}

func FromProducer(producer Producer) *Props {
	return &Props{
		producer:        producer,
		spawner:         defaultSpawner,
		mailboxProducer: defaultMailboxProducer,
	}
}
