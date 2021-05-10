package message

import (
	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/conn/packet"
	"google.golang.org/protobuf/proto"
)

//同步消息
func UnpackReq(p *packet.Packet) *rpc.Request {
	var res = &rpc.Request{}
	proto.Unmarshal(p.Body, res)
	return res
}

func NewRequst(className, methodName string, args ...RequestOption) *rpc.Request {
	var req = &rpc.Request{
		ServerName: className,
		MethodName: methodName,
	}

	for _, f := range args {
		f(req)
	}

	return req
}

//可选参数
type RequestOption func(*rpc.Request)

func SetReqData(data proto.Message) RequestOption {
	return func(p *rpc.Request) {
		p.Body, _ = proto.Marshal(data)
	}
}

func SetReqMetadata(metadata map[string]string) RequestOption {
	return func(p *rpc.Request) {
		p.Metadata = metadata
	}
}
