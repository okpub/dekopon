package actor

import (
	"fmt"
)

type ActorOptions struct {
	//context.Context
	ID   int
	Name string
	Addr string
}

var defaultActorOptions = &ActorOptions{}

func NewOptions(args []ActorOption) *ActorOptions {
	var options = defaultActorOptions.Options()
	return options.Filler(args)
}

func (options *ActorOptions) Filler(args []ActorOption) *ActorOptions {
	for _, f := range args {
		f(options)
	}
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

func SetID(id int) ActorOption {
	return func(p *ActorOptions) {
		p.ID = id
	}
}

func SetName(name string) ActorOption {
	return func(p *ActorOptions) {
		p.Name = name
	}
}
