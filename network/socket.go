package network

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/skimmer/actor"
	"github.com/skimmer/conn/codec"
	"github.com/skimmer/conn/packet"
	"github.com/skimmer/mailbox"
)

/*
 * 问答socket
 */
type Socket struct {
	mailbox.TaskBuffer
	mailbox.InvokerMessage
	SocketOptions

	packetChan mailbox.TaskBuffer
}

func NewSocket(args ...SocketOption) *Socket {
	var options = NewOptions()
	options.Filler(args)
	return &Socket{
		SocketOptions: options,
		TaskBuffer:    mailbox.MakeBuffer(options.PendingNum),
		packetChan:    mailbox.MakeBlock(),
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

//代替Start
func (socket *Socket) Serve(ctx context.Context, conn net.Conn, err error) error {
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
func (socket *Socket) Start(ctx context.Context) error {
	panic(fmt.Errorf("please apply Serve with Error"))
}

func (socket *Socket) run(ctx context.Context, conn net.Conn) (err error) {
	socket.InvokeSystemMessage(EVENT_OPEN)
	//异步写入
	go socket.ListenAndWrite(conn)
	//同步读取
	go func() {
		defer socket.packetChan.Close()
		socket.ListenAndRead(conn)
	}()

	for data := range socket.packetChan {
		socket.InvokeUserMessage(data)
	}
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
		if packets, err = socket.readPackets(buf, conn); err == nil {
			for _, message := range packets {
				//socket.InvokeUserMessage(message)
				socket.packetChan <- message
			}
		} else {
			if Temporary(err) {
				socket.packetChan <- &TempErr{Err: err}
				//socket.InvokeSystemMessage(&TempErr{Err: err})
			} else {
				break
			}

		}
	}
	return
}

func (socket *Socket) readPackets(buf *bytes.Buffer, conn net.Conn) (packets []*packet.Packet, err error) {
	var (
		n              int
		totalProcessed int
		body           [ReadBuffSize]byte
	)

	if n, err = conn.Read(body[0:]); n > 0 {
		buf.Write(body[:n])
	}

	if err == nil {
		packets, err = socket.Decoder.Decode(buf.Bytes())
		for _, message := range packets {
			totalProcessed += codec.HeadSize + message.Len()
		}
		buf.Next(totalProcessed)
	} else {
		fmt.Println(err.Error())
	}
	return
}
