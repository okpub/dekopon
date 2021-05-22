package actor

import (
	"github.com/okpub/dekopon/mailbox"
	"github.com/okpub/dekopon/utils"
)

var (
	defaultMailboxProducer = mailbox.Unbounded()
	defaultDispatcher      = mailbox.NewDefaultDispatcher()
	defaultSpawner         = func(parent SpawnContext, props *Props, options *ActorOptions) PID {
		var (
			done    = utils.MakeDone()
			mbox    = props.NewMailbox()
			pid     = WithPID(NewDfaultProcess(done.Done(), mbox), options)
			invoker = NewContext(parent.ChildOf(pid), props)
		)

		mbox.RegisterHander(invoker, props.GetDispatcher())

		go func() {
			defer done.Shutdown()
			invoker.InvokeSystemMessage(EVENT_START)
			mbox.Start()
			invoker.InvokeSystemMessage(EVENT_STOP)
		}()

		return pid
	}
)

//目前够用
type Props struct {
	producer        Producer
	spawner         SpawnFunc
	mailboxProducer mailbox.Producer
	dispatcher      mailbox.Dispatcher

	//spawn
	spawnMiddleware      []SpawnMiddleware
	spawnMiddlewareChain SpawnFunc

	//middkeware
	receiverMiddleware      []ReceiverMiddleware
	receiverMiddlewareChain ReceiverFunc
	senderMiddleware        []SenderMiddleware
	senderMiddlewareChain   SenderFunc
	contextDecorator        []ContextDecorator
	contextDecoratorChain   ContextDecoratorFunc
}

func (props *Props) NewActor() Actor {
	return props.producer()
}

func (props *Props) NewMailbox() mailbox.Mailbox {
	if props.mailboxProducer == nil {
		return defaultMailboxProducer()
	}
	return props.mailboxProducer()
}

func (props *Props) GetDispatcher() mailbox.Dispatcher {
	if props.dispatcher == nil {
		return defaultDispatcher
	}
	return props.dispatcher
}

func (props Props) Options() Props {
	return props
}

//context apply
func (props *Props) spawn(parent SpawnContext, options *ActorOptions) (pid PID) {
	if chain := props.spawnMiddlewareChain; chain != nil {
		return chain(parent, props, options)
	}
	if props.spawner == nil {
		return defaultSpawner(parent, props, options)
	}
	return props.spawner(parent, props, options)
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
		producer: producer,
	}
}
