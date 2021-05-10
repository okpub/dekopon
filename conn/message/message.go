package message

import (
	"github.com/okpub/dekopon/bean/message/common"
	"github.com/okpub/dekopon/conn/packet"
	"google.golang.org/protobuf/proto"
)

const (
	defaultVersion    = 100
	defaultSubVersion = 1
)

type DataMessage interface {
	GetBody() []byte
}

//unpack
func UnPack(p *packet.Packet) *common.UserMessage {
	var res = &common.UserMessage{}
	proto.Unmarshal(p.Body, res)
	return res
}

func Pack(cmd int32, args ...MessageOption) *common.UserMessage {
	var req = &common.UserMessage{
		Header: &common.MessageHeader{
			Cmd:        cmd,
			Version:    defaultVersion,
			SubVersion: defaultSubVersion,
		},
	}

	for _, o := range args {
		o(req)
	}

	return req
}

func GetMessage(m DataMessage, any proto.Message) (err error) {
	err = proto.Unmarshal(m.GetBody(), any)
	return
}

//消息头可选参数
type MessageOption func(*common.UserMessage)

//消息类型
func SetMessageType(n int32) MessageOption {
	return func(p *common.UserMessage) {
		p.Header.MessageType = n
	}
}

//服务器id
func SetServer(n int32) MessageOption {
	return func(p *common.UserMessage) {
		p.Header.ServerId = n
	}
}

func SetMessageData(any proto.Message) MessageOption {
	return func(p *common.UserMessage) {
		p.Body, _ = proto.Marshal(any)
	}
}

//可选参数
func SetSession(session *common.Session) MessageOption {
	return func(p *common.UserMessage) {
		p.Session = session
	}
}

func SetSessionMetaData(metadata map[string]string) MessageOption {
	return func(p *common.UserMessage) {
		if p.Session == nil {
			p.Session = &common.Session{}
		}
		p.Session.MetaData = metadata
	}
}

func SetSessionHeader(header *common.SessionHeader) MessageOption {
	return func(p *common.UserMessage) {
		if p.Session == nil {
			p.Session = &common.Session{}
		}
		p.Session.Header = header
	}
}
