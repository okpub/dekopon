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

func (client *Client) Start() {
	go client.Serve(context.Background())
}

func (client *Client) Serve(ctx context.Context) {
	var (
		conn, err     = client.Connect()
		child, cancel = context.WithCancel(ctx)
	)
	defer cancel()
	if err == nil {
		go func() {
			client.ListenAndWrite(child, conn)
		}()

		client.ListenAndRead(conn)
	} else {
		client.Close()
	}
}

func (client *Client) ListenAndWrite(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	var (
		sendCh        = client.TaskBuffer
		handleRequest = func(event actor.Request) {
			if body, err := client.Encoder.Encode(event.Message()); err == nil {
				if _, err := conn.Write(body); err == nil {
					select {
					case <-ctx.Done():
						//cancel error
					case p := <-client.packetChan:
						event.Respond(p)
					}
				}
			} else {
				fmt.Println(err)
			}
		}
	)

	for message := range sendCh {
		switch event := message.(type) {
		case actor.Request:
			handleRequest(event)
		default:
			fmt.Println("can't handler other message type:", message)
		}
	}
}

func (client *Client) ListenAndRead(conn net.Conn) (err error) {
	var (
		sendChan = client.TaskBuffer
		buf      = bytes.NewBuffer(nil)
		packets  []*packet.Packet
	)
	defer sendChan.Close()
	for {
		if packets, err = network.ReadPackets(buf, conn, client.Decoder); err == nil {
			for _, message := range packets {
				client.packetChan <- message
			}
		} else {
			break
		}
	}
	return
}

//request func
func (client *Client) Request(req *rpc.Request) (*rpc.Response, error) {
	return client.RequestCtx(context.Background(), req)
}

func (client *Client) RequestCtx(ctx context.Context, req *rpc.Request) (res *rpc.Response, err error) {
	var data interface{}
	if data, err = client.CallUserMessage(ctx, req); err == nil {
		res, err = message.UnpackRes(data)
	}
	return
}
