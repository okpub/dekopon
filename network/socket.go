package network

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/conn/codec"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/mailbox"
)

/*
 * 问答socket
 */
type Socket struct {
	mailbox.TaskBuffer
	mailbox.InvokerMessage
	*SocketOptions
}

func NewSocket(args ...SocketOption) *Socket {
	var options = NewOptions()
	return &Socket{
		SocketOptions: options.Filler(args),
		TaskBuffer:    options.NewChannel(),
	}
}

func WithSocket(conn net.Conn, args ...SocketOption) *Socket {
	var socket = NewSocket(args...)
	socket.WithAddr(conn.RemoteAddr().String())
	return socket
}

func (socket *Socket) RegisterHander(invoker mailbox.InvokerMessage) {
	socket.InvokerMessage = invoker
}

//不需要连接
func (socket *Socket) ServeConn(ctx context.Context, conn net.Conn) (err error) {
	defer socket.InvokeSystemMessage(actor.EVENT_STOP)
	defer socket.InvokeSystemMessage(EVENT_CLOSED)

	socket.InvokeSystemMessage(actor.EVENT_START)
	err = socket.run(ctx, conn)
	return
}

//override public
func (socket *Socket) Start(ctx context.Context) error {
	var (
		conn, err = socket.Connect()
	)

	defer socket.InvokeSystemMessage(actor.EVENT_STOP)
	defer socket.InvokeSystemMessage(EVENT_CLOSED)

	socket.InvokeSystemMessage(actor.EVENT_START)
	if err == nil {
		err = socket.run(ctx, conn)
	} else {
		socket.InvokeSystemMessage(&DialError{Err: err})
		socket.Close()
	}
	return err
}

func (socket *Socket) run(ctx context.Context, conn net.Conn) (err error) {
	socket.InvokeSystemMessage(EVENT_OPEN)
	//异步写入
	go socket.ListenAndWrite(conn)
	//同步读取
	socket.ListenAndRead(conn)
	return
}

/*
* 监听写入
 */
func (socket *Socket) ListenAndWrite(conn net.Conn) (err error) {
	var (
		sendCh = socket.TaskBuffer
		body   []byte
	)
	defer conn.Close()
	for message := range sendCh {
		if body, err = socket.Encoder.Encode(actor.GetMessage(message)); err == nil {
			_, err = conn.Write(body)
		}
	}
	return
}

/*
* 监听读取
 */
func (socket *Socket) ListenAndRead(conn net.Conn) (err error) {
	var (
		sendCh  = socket.TaskBuffer
		buf     = bytes.NewBuffer(nil)
		packets []*packet.Packet
	)
	defer sendCh.Close()
	for {
		if packets, err = ReadPackets(buf, conn, socket.Decoder); err == nil {
			for _, message := range packets {
				socket.InvokeUserMessage(message)
			}
		} else {
			if Temporary(err) {
				socket.InvokeSystemMessage(&TempErr{Err: err})
			} else {
				break
			}

		}
	}
	return
}

func ReadPackets(buf *bytes.Buffer, conn net.Conn, decoder codec.PacketDecoder) (packets []*packet.Packet, err error) {
	var (
		n              int
		totalProcessed int
		body           [ReadBuffSize]byte
	)

	if n, err = conn.Read(body[0:]); n > 0 {
		buf.Write(body[:n])
	}

	if err == nil {
		packets, err = decoder.Decode(buf.Bytes())
		for _, message := range packets {
			totalProcessed += codec.HeadSize + message.Len()
		}
		buf.Next(totalProcessed)
	} else {
		fmt.Println(err.Error())
	}
	return
}
