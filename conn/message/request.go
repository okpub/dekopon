package message

import (
	"fmt"

	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/conn/packet"
	"google.golang.org/protobuf/proto"
)

func Request(p *packet.Packet) *rpc.Request {
	var res = &rpc.Request{}
	proto.Unmarshal(p.Body, res)
	return res
}

/*
* 得到回执
 */
func Ask(data interface{}) (res *rpc.Response, err error) {
	switch p := data.(type) {
	case *packet.Packet:
		res = &rpc.Response{}
		err = proto.Unmarshal(p.Body, res)
	case *rpc.Response:
		res = p
	default:
		err = fmt.Errorf("the res can't unpack %T", data)
	}
	return
}

func ResOk(args ...ResponseOption) *rpc.Response {
	var res = &rpc.Response{ErrMsg: "request ok!"}
	for _, f := range args {
		f(res)
	}
	return res
}

func ResErr(code int32, args ...ResponseOption) *rpc.Response {
	var res = &rpc.Response{ErrCode: code, ErrMsg: "request err!"}
	for _, f := range args {
		f(res)
	}
	return res
}

//options
type ResponseOption func(*rpc.Response)

func SetResponseData(data proto.Message) ResponseOption {
	return func(p *rpc.Response) {
		p.Body, _ = proto.Marshal(data)
	}
}

func SetResponseErr(errmsg string) ResponseOption {
	return func(p *rpc.Response) {
		p.ErrMsg = errmsg
	}
}
