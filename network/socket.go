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
	"github.com/okpub/dekopon/utils"
)

/*
 * 问答socket
 */
type SocketProcess struct {
	mailbox.TaskBuffer
	mailbox.InvokerMessage
	SocketOptions
}

func NewSocket(options SocketOptions) *SocketProcess {
	return &SocketProcess{SocketOptions: options, TaskBuffer: mailbox.MakeBuffer(options.PendingNum)}
}

func WithSocket(options SocketOptions, conn net.Conn) *SocketProcess {
	var socket = NewSocket(options)
	socket.WithAddr(conn.RemoteAddr().String()).WithNetwork(WEB)
	return socket
}

func (socket *SocketProcess) RegisterHander(invoker mailbox.InvokerMessage) {
	socket.InvokerMessage = invoker
}

//代替Start
func (socket *SocketProcess) Serve(ctx context.Context, conn net.Conn, err error) error {
	defer socket.InvokeSystemMessage(actor.EVENT_STOP)
	defer socket.InvokeSystemMessage(EVENT_CLOSED)

	socket.InvokeSystemMessage(actor.EVENT_START)
	if err == nil {
		err = socket.run(ctx, conn)
	} else {
		socket.InvokeSystemMessage(&DialError{Err: err})
		//close buffer
		socket.Close()
	}
	return err
}

//override public
func (socket *SocketProcess) Start(ctx context.Context) error {
	panic(fmt.Errorf("please apply Serve with Error"))
}

func (socket *SocketProcess) run(ctx context.Context, conn net.Conn) (err error) {
	var (
		wg actor.WaitGroup
	)
	socket.InvokeSystemMessage(EVENT_OPEN)
	//异步写入
	wg.Wrap(func() {
		socket.ListenAndWrite(conn)
	})
	//同步读取
	err = socket.ListenAndRead(conn)
	wg.Wait()
	return
}

/*
* 监听写入
 */
func (socket *SocketProcess) ListenAndWrite(conn net.Conn) (err error) {
	var (
		sendCh = socket.TaskBuffer
		body   []byte
	)
	defer conn.Close()
	for message := range sendCh {
		if body, err = socket.Encoder.Encode(message); err == nil {
			_, err = conn.Write(body)
		}
	}
	return
}

/*
* 监听读取
 */
func (socket *SocketProcess) ListenAndRead(conn net.Conn) (err error) {
	var (
		sendCh  = socket.TaskBuffer
		buf     = bytes.NewBuffer(nil)
		packets []*packet.Packet
	)
	defer sendCh.Close()
	for {
		packets, err = socket.readPackets(buf, conn)
		for _, message := range packets {
			socket.InvokeUserMessage(message)
		}

		if utils.Die(err) {
			if Temporary(err) {
				socket.InvokeSystemMessage(&TempErr{Err: err})
			} else {
				break
			}
		}
	}
	return
}

func (socket *SocketProcess) readPackets(buf *bytes.Buffer, conn net.Conn) (packets []*packet.Packet, err error) {
	var (
		n              int
		totalProcessed int
		body           [ReadBuffSize]byte
	)

	if n, err = conn.Read(body[0:]); n > 0 {
		_, err = buf.Write(body[:n])
	}

	packets, err = socket.Decoder.Decode(buf.Bytes())
	for _, message := range packets {
		totalProcessed += codec.HeadSize + message.Len()
	}
	buf.Next(totalProcessed)
	return
}
