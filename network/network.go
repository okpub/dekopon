package network

import (
	"context"
	"net"
	"time"
)

type Handler func(context.Context, net.Conn)

const (
	TemporaryInterval = time.Millisecond * 100 //临时报错间隔时间
	ReadBuffSize      = 1024 * 2               //默认读取缓存区大小
	LimitConnsNone    = -1                     //无限server连接
	ShakeHandsTime    = time.Second * 3        //握手时间
	PingTime          = time.Second * 60       //心跳时间
	DialTimeout       = time.Second * 8        //连接超时
	ReadTimeout       = time.Second * 3        //第一次读超时
	WritePendingNum   = 20                     //写入队列长度
)

type Server interface {
	Options() ServerOptions
	ListenAndServe(context.Context, Handler) error
	Close() error
}

/*
* 临时错误
 */
func Temporary(err error) bool {
	var temp, ok = err.(net.Error)
	return ok && temp.Temporary()
}
