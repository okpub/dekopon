package client

import (
	"bytes"
	"context"
	"net"

	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/common"
	"github.com/okpub/dekopon/conn/codec"
	"github.com/okpub/dekopon/conn/packet"
	"github.com/okpub/dekopon/mailbox"
	"github.com/okpub/dekopon/network"
	"github.com/okpub/dekopon/utils"
)

type Client struct {
	network.SocketOptions
	requestChan chan actor.Request
	packetChan  chan *packet.Packet
}

func (client *Client) handleWrite(conn net.Conn) {
	for data := range client.requestChan {
		if body, err := client.Encoder.Encode(data.Message()); err == nil {
			_, err = conn.Write(body)
			if err == nil {
				data.Respond(<-client.packetChan)
			}
		}
	}
}

func (client *Client) Start(conn net.Conn) (err error) {
	go client.handleWrite(conn)
	go client.handlePackets(conn)
	return
}

func (client *Client) handleWrite(conn net.Conn) (err error) {
	var body []byte
	for message := range client.sendChan {
		if body, err = client.Encoder.Encode(message); err == nil {
			_, err = conn.Write(body)
		}
	}
	return
}

func (client *Client) handlePackets(conn net.Conn) (err error) {
	var (
		buf     = bytes.NewBuffer(nil)
		packets []*packet.Packet
	)
	for {
		if packets, err = client.readPackets(buf, conn); utils.Die(err) {
			break
		}
		for _, p := range packets {
			client.packetChan <- p
		}
	}
	return
}

func (client *Client) readPackets(buf *bytes.Buffer, conn net.Conn) (packets []*packet.Packet, err error) {
	var (
		body [network.ReadBuffSize]byte
		n    int
	)

	if n, err = conn.Read(body[0:]); n > 0 {
		buf.Write(body[:n])
	}

	client.Decoder.Decode(buf.Bytes())

	totalProcessed := 0
	for _, p := range packets {
		totalProcessed += codec.HeadSize + p.Len()
	}

	buf.Next(totalProcessed)
	return
}

func (client *Client) SendRequest(req *common.Request) (res *common.Response, err error) {
	var data = mailbox.NewRequest(req)
	client.requestChan <- data

	data.Body(context.Background())
	return
}
