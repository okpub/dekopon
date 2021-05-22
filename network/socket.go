package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/conn/codec"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/mailbox"
	"github.com/okpub/dekopon/utils"
)

/*
 * 问答socket
 */
type Socket struct {
	utils.TaskBuffer
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

func (socket *Socket) RegisterHander(invoker mailbox.InvokerMessage, distacher mailbox.Dispatcher) {
	socket.InvokerMessage = invoker
}

//不需要连接
func (socket *Socket) ServeConn(conn net.Conn) (err error) {
	socket.InvokeSystemMessage(EVENT_OPEN)
	//异步写入
	go socket.ListenAndWrite(conn)
	//同步读取
	err = socket.ListenAndRead(conn)
	return
}

//override public
func (socket *Socket) Start() {
	var (
		conn, err = socket.Connect()
	)
	if err == nil {
		socket.ServeConn(conn)
	} else {
		socket.InvokeSystemMessage(&DialError{Err: err})
		socket.Close()
	}
}

/*
* 监听写入
 */
func (socket *Socket) ListenAndWrite(conn net.Conn) (err error) {
	var (
		sendChan = socket.TaskBuffer
		body     []byte
	)
	defer conn.Close()
	for message := range sendChan {
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
		sendChan = socket.TaskBuffer
		buf      = bytes.NewBuffer(nil)
		packets  []*packet.Packet
	)
	defer utils.SafeClose(sendChan)
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

//mailbox
func (socket *Socket) CallUserMessage(ctx context.Context, message interface{}) (res interface{}, err error) {
	panic(errors.New("不支持"))
}

func (socket *Socket) PostUserMessage(ctx context.Context, message interface{}) error {
	return utils.SendCtx(socket.TaskBuffer, ctx, message)
}

func (socket *Socket) PostSystemMessage(message interface{}) error {
	panic(errors.New("不支持"))
}
