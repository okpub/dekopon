package actor

import (
	"context"
)

func (props *Props) WithSpawnFunc(spawn SpawnFunc) *Props {
	props.spawner = spawn
	return props
}

func (props *Props) WithSpawnMiddleware(middleware ...SpawnMiddleware) *Props {
	props.spawnMiddleware = append(props.spawnMiddleware, middleware...)

	props.spawnMiddlewareChain = makeSpawnMiddlewareChain(props.spawnMiddleware,
		func(parent SpawnContext, _ *Props, options *ActorOptions) PID {
			return props.spawn(parent, options)
		})

	return props
}

//recv
func (props *Props) WithReceiverMiddleware(middleware ...ReceiverMiddleware) *Props {
	props.receiverMiddleware = append(props.receiverMiddleware, middleware...)

	props.receiverMiddlewareChain = makeReceiverMiddlewareChain(props.receiverMiddleware,
		func(ctx ReceiverContext, envelope MessageEnvelope) {
			ctx.Received(envelope)
		})

	return props
}

//send
func (props *Props) WithSenderMiddleware(middleware ...SenderMiddleware) *Props {
	props.senderMiddleware = append(props.senderMiddleware, middleware...)

	props.senderMiddlewareChain = makeSenderMiddlewareChain(props.senderMiddleware,
		func(ctx context.Context, sender SenderContext, pid PID, envelope MessageEnvelope) error {
			return pid.sendUserMessage(ctx, envelope)
		})

	return props
}

//extends
func (props *Props) WithContextDecorator(contextDecorator ...ContextDecorator) *Props {
	props.contextDecorator = append(props.contextDecorator, contextDecorator...)

	props.contextDecoratorChain = makeContextDecoratorChain(props.contextDecorator,
		func(ctx ActorContext) ActorContext {
			return ctx
		})

	return props
}
