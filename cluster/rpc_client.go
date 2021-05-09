package cluster

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/skimmer/actor"
	"github.com/skimmer/bean/message/rpc"
	"github.com/skimmer/conn/codec"
	"github.com/skimmer/conn/message"
	"github.com/skimmer/conn/packet"
	"github.com/skimmer/mailbox"
	"github.com/skimmer/network"
)

//应答客户端
type Client struct {
	*network.SocketOptions
	mailbox.TaskBuffer
	packetChan chan *packet.Packet
}

func NewClient(args ...network.SocketOption) *Client {
	var options = network.NewOptions()
	return &Client{
		SocketOptions: options.Filler(args),
		TaskBuffer:    mailbox.MakeBuffer(options.PendingNum),
		packetChan:    make(chan *packet.Packet)}
}

func (client *Client) Start(ctx context.Context) {
	go client.Serve(ctx)
}

func (client *Client) Serve(parent context.Context) {
	var (
		conn, err   = client.Connect()
		ctx, cancel = context.WithCancel(parent)
	)
	defer cancel()
	defer client.Close()
	if err == nil {
		go client.ListenAndWrite(ctx, conn)

		client.ListenAndRead(ctx, conn)
	}
}

func (client *Client) ListenAndWrite(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for message := range client.TaskBuffer {
		switch event := message.(type) {
		case actor.Request:
			client.processRequest(ctx, conn, event)
		default:
			fmt.Println("can't handler message type:", message)
		}
	}
}

func (client *Client) processRequest(ctx context.Context, conn net.Conn, event actor.Request) {
	if body, err := client.Encoder.Encode(event.Message()); err == nil {
		if _, err := conn.Write(body); err == nil {
			select {
			case <-ctx.Done():
				//cancel error
			case p := <-client.packetChan:
				event.Respond(p)
			}
		}
	}
}

func (client *Client) Request(ctx context.Context, req *rpc.Request) (res *rpc.Response, err error) {
	var data interface{}
	data, err = client.CallUserMessage(ctx, req)
	if err == nil {
		res, err = message.Ask(data.(*packet.Packet))
	}
	return
}

/*
* 监听读取
 */
func (client *Client) ListenAndRead(ctx context.Context, conn net.Conn) (err error) {
	var (
		buf     = bytes.NewBuffer(nil)
		packets []*packet.Packet
	)
	for {
		if packets, err = client.readPackets(buf, conn); err == nil {
			for _, message := range packets {
				select {
				case <-ctx.Done():
					//cancel message
				case client.packetChan <- message:
					//read packet
				}
			}
		} else {
			//忽略临时报错, 直接关闭
			break
		}
	}
	return
}

func (client *Client) readPackets(buf *bytes.Buffer, conn net.Conn) (packets []*packet.Packet, err error) {
	var (
		n              int
		totalProcessed int
		body           [network.ReadBuffSize]byte
	)

	if n, err = conn.Read(body[0:]); n > 0 {
		buf.Write(body[:n])
	}

	if err == nil {
		packets, err = client.Decoder.Decode(buf.Bytes())
		for _, message := range packets {
			totalProcessed += codec.HeadSize + message.Len()
		}
		buf.Next(totalProcessed)
	} else {
		fmt.Println(err.Error())
	}
	return
}
