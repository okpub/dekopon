package actor

import (
	"testing"
)

func middleware(called *int) ReceiverMiddleware {
	return func(next ReceiverFunc) ReceiverFunc {
		return func(ctx ReceiverContext, env MessageEnvelope) {
			//env.Message = env.Message.(int) + 1
			//*called = env.Message.(int)

			next(ctx, env)
		}
	}
}

func TestMakeReceiverMiddleware_CallsInCorrectOrder(t *testing.T) {
	var c [3]int

	r := []ReceiverMiddleware{
		middleware(&c[0]),
		middleware(&c[1]),
		middleware(&c[2]),
	}

	env := MSG(10)

	chain := makeReceiverMiddlewareChain(r, func(_ ReceiverContext, env MessageEnvelope) {})
	chain(nil, env)

}

func TestMakeInboundMiddleware_ReturnsNil(t *testing.T) {
	makeReceiverMiddlewareChain([]ReceiverMiddleware{}, func(_ ReceiverContext, _ MessageEnvelope) {})
}
