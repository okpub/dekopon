package cluster

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/conn/message"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/mailbox"
	"github.com/okpub/dekopon/network"
)

//同步应答客户端
type Client struct {
	*network.SocketOptions
	mailbox.TaskBuffer
	packetChan chan *packet.Packet
}

func NewClient(args ...network.SocketOption) *Client {
	var options = network.NewOptions()
	return &Client{
		SocketOptions: options.Filler(args),
		TaskBuffer:    options.NewChannel(),
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

		client.ListenAndRead(conn)
	}
}

func (client *Client) ListenAndWrite(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	var sendCh = client.TaskBuffer
	for message := range sendCh {
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
func (client *Client) ListenAndRead(conn net.Conn) (err error) {
	var (
		buf     = bytes.NewBuffer(nil)
		packets []*packet.Packet
	)
	for {
		if packets, err = network.ReadPackets(buf, conn, client.Decoder); err == nil {
			for _, message := range packets {
				client.packetChan <- message
			}
		} else {
			//忽略临时报错, 直接关闭
			break
		}
	}
	return
}
