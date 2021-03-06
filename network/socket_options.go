package network

import (
	"net"
	"time"

	"github.com/okpub/dekopon/conn/codec"
	"github.com/okpub/dekopon/utils"
)

const (
	TCP = "tcp"
	WEB = "wss"
)

var (
	defaultSocketOptions = SocketOptions{
		Network:    TCP,
		PendingNum: WritePendingNum,
		DialDur:    DialTimeout,
		Encoder:    codec.NewPomeloPacketEncoder(),
		Decoder:    codec.NewPomeloPacketDecoder(),
	}
)

func NewOptions(args ...SocketOption) SocketOptions {
	var options = defaultSocketOptions
	options.Filler(args)
	return options
}

type SocketOptions struct {
	Addr       string        //远端(需要用)
	Network    string        //类型
	PendingNum int           //缓冲大小
	DialDur    time.Duration //连接超时

	//option
	SendDur time.Duration //发送超时
	ReadDur time.Duration //读包时间

	//codec
	Encoder codec.PacketEncoder
	Decoder codec.PacketDecoder
}

func (options *SocketOptions) Filler(args []SocketOption) *SocketOptions {
	for _, f := range args {
		f(options)
	}
	return options
}

func (options *SocketOptions) Connect() (net.Conn, error) {
	return Dial(options)
}

func (options *SocketOptions) WithAddr(addr string) *SocketOptions {
	options.Addr = addr
	return options
}

func (options *SocketOptions) WithNetwork(kind string) *SocketOptions {
	options.Network = kind
	return options
}

func (options *SocketOptions) NewChannel() utils.TaskBuffer {
	return utils.MakeBuffer(options.PendingNum)
}

func (options SocketOptions) Options() *SocketOptions {
	return &options
}

//socket可选参数
type SocketOption func(*SocketOptions)

func SetDialAddr(addr string) SocketOption {
	return func(p *SocketOptions) {
		p.Addr = addr
	}
}

func SetDialNetwork(addr string) SocketOption {
	return func(p *SocketOptions) {
		p.Network = addr
	}
}
