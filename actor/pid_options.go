package actor

import (
	"fmt"
)

type ActorOptions struct {
	//context.Context
	ID   int32
	Name string
	Addr string
}

var defaultActorOptions = &ActorOptions{}

func NewOptions(args []ActorOption) *ActorOptions {
	var options = defaultActorOptions.Options()
	return options.Filler(args)
}

func WithOptions(args ...ActorOption) *ActorOptions {
	var options = defaultActorOptions.Options()
	return options.Filler(args)
}

func (options *ActorOptions) Filler(args []ActorOption) *ActorOptions {
	for _, f := range args {
		f(options)
	}
	return options
}

func (options *ActorOptions) WithAddr(addr string) *ActorOptions {
	options.Addr = addr
	return options
}

func (options *ActorOptions) WithName(name string) *ActorOptions {
	options.Name = name
	return options
}

func (options *ActorOptions) WithId(id int32) *ActorOptions {
	options.ID = id
	return options
}

func (options ActorOptions) Options() ActorOptions {
	return options
}

func (options *ActorOptions) String() string {
	return fmt.Sprintf("id:%d name:%s addr:%s", options.ID, options.Name, options.Addr)
}

//可选参数
type ActorOption func(*ActorOptions)

func SetID(id int32) ActorOption {
	return func(p *ActorOptions) {
		p.ID = id
	}
}

func SetName(name string) ActorOption {
	return func(p *ActorOptions) {
		p.Name = name
	}
}
