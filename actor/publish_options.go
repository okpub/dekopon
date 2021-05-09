package actor

import (
	"context"
	"time"
)

/*
 * 发布参数(MessageEnvelope)
 */
type PublishOptions struct {
	context.Context
	Message interface{}
}

func NewPublish(message interface{}) PublishOptions {
	return PublishOptions{
		Message: message,
		Context: context.Background(),
	}
}

func (options *PublishOptions) Filler(args []PublishOption) *PublishOptions {
	for _, f := range args {
		f(options)
	}
	return options
}

//可选参数(Message中途可能被劫持而改变)
type PublishOption func(*PublishOptions)

func SetContext(ctx context.Context) PublishOption {
	return func(p *PublishOptions) {
		p.Context = ctx
	}
}

func SetTimeout(dur time.Duration) PublishOption {
	return func(p *PublishOptions) {
		p.Context, _ = context.WithTimeout(context.Background(), dur)
	}
}
