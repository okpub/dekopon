package actor

import (
	"fmt"
)

type PIDOptions struct {
	ID   int
	Name string
	Addr string
}

func NewOptions(args []PIDOption) *PIDOptions {
	var options = PIDOptions{}
	return options.Filler(args)
}

func (options *PIDOptions) Filler(args []PIDOption) *PIDOptions {
	for _, f := range args {
		f(options)
	}
	return options
}

func (options PIDOptions) Options() PIDOptions {
	return options
}

func (options *PIDOptions) String() string {
	return fmt.Sprintf("id:%d name:%s addr:%s", options.ID, options.Name, options.Addr)
}

//可选参数
type PIDOption func(*PIDOptions)

func SetID(id int) PIDOption {
	return func(p *PIDOptions) {
		p.ID = id
	}
}

func SetName(name string) PIDOption {
	return func(p *PIDOptions) {
		p.Name = name
	}
}

func SetAddr(addr string) PIDOption {
	return func(p *PIDOptions) {
		p.Addr = addr
	}
}
