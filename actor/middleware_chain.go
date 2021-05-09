package actor

import (
	"context"
)

type ValueFunc func(context.Context) context.Context
type ValueMiddleware func(next ValueFunc) ValueFunc

//middleware
type ReceiverFunc func(ReceiverContext, MessageEnvelope)
type ReceiverMiddleware func(next ReceiverFunc) ReceiverFunc

type SenderFunc func(context.Context, SenderContext, PID, MessageEnvelope) error
type SenderMiddleware func(next SenderFunc) SenderFunc

type ContextDecoratorFunc func(ActorContext) ActorContext
type ContextDecorator func(next ContextDecoratorFunc) ContextDecoratorFunc

type SpawnFunc func(SpawnContext, *Props, *PIDOptions) PID
type SpawnMiddleware func(next SpawnFunc) SpawnFunc

func makeReceiverMiddlewareChain(receiverMiddleware []ReceiverMiddleware, lastReceiver ReceiverFunc) ReceiverFunc {
	if len(receiverMiddleware) == 0 {
		return nil
	}

	h := receiverMiddleware[len(receiverMiddleware)-1](lastReceiver)
	for i := len(receiverMiddleware) - 2; i >= 0; i-- {
		h = receiverMiddleware[i](h)
	}
	return h
}

func makeSenderMiddlewareChain(senderMiddleware []SenderMiddleware, lastSender SenderFunc) SenderFunc {
	if len(senderMiddleware) == 0 {
		return nil
	}

	h := senderMiddleware[len(senderMiddleware)-1](lastSender)
	for i := len(senderMiddleware) - 2; i >= 0; i-- {
		h = senderMiddleware[i](h)
	}
	return h
}

func makeContextDecoratorChain(decorator []ContextDecorator, lastDecorator ContextDecoratorFunc) ContextDecoratorFunc {
	if len(decorator) == 0 {
		return nil
	}

	h := decorator[len(decorator)-1](lastDecorator)
	for i := len(decorator) - 2; i >= 0; i-- {
		h = decorator[i](h)
	}
	return h
}

func makeSpawnMiddlewareChain(spawnMiddleware []SpawnMiddleware, lastSpawn SpawnFunc) SpawnFunc {
	if len(spawnMiddleware) == 0 {
		return nil
	}

	h := spawnMiddleware[len(spawnMiddleware)-1](lastSpawn)
	for i := len(spawnMiddleware) - 2; i >= 0; i-- {
		h = spawnMiddleware[i](h)
	}
	return h
}

func makeValueMiddlewareChain(valueMiddleware []ValueMiddleware, lastSpawn ValueFunc) ValueFunc {
	if len(valueMiddleware) == 0 {
		return nil
	}

	h := valueMiddleware[len(valueMiddleware)-1](lastSpawn)
	for i := len(valueMiddleware) - 2; i >= 0; i-- {
		h = valueMiddleware[i](h)
	}
	return h
}
