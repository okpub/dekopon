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
		func(parent SpawnContext, _ *Props, options *PIDOptions) PID {
			return props.spawn(parent, options)
		})

	return props
}

//values
func (props *Props) WithValue(key, value interface{}) *Props {
	var middleware = func(next ValueFunc) ValueFunc {
		return func(ctx context.Context) context.Context {
			return next(context.WithValue(ctx, key, value))
		}
	}
	return props.WithValueMiddleware(middleware)
}

func (props *Props) WithValueMiddleware(middleware ...ValueMiddleware) *Props {
	props.valueMiddleware = append(props.valueMiddleware, middleware...)

	props.valueMiddlewareChain = makeValueMiddlewareChain(props.valueMiddleware,
		func(ctx context.Context) context.Context {
			return ctx
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
