package message

import (
	"fmt"

	"github.com/okpub/dekopon/bean/message/rpc"
	"github.com/okpub/dekopon/conn/packet"
	"google.golang.org/protobuf/proto"
)

var (
	resOkString  = "request ok!"
	resErrString = "request err!"
)

func UnpackRes(data interface{}, failed error) (res *rpc.Response, err error) {
	if err = failed; err == nil {
		switch p := data.(type) {
		case *packet.Packet:
			res = &rpc.Response{}
			err = proto.Unmarshal(p.Body, res)
		case *rpc.Response:
			res = p
		default:
			err = fmt.Errorf("the res can't unpack %T", data)
		}
	}
	return
}

func ResOk(args ...ResponseOption) *rpc.Response {
	var res = &rpc.Response{ErrMsg: resOkString}
	for _, f := range args {
		f(res)
	}
	return res
}

func ResErr(code int32, args ...ResponseOption) *rpc.Response {
	var res = &rpc.Response{ErrCode: code, ErrMsg: resErrString}
	for _, f := range args {
		f(res)
	}
	return res
}

//options
type ResponseOption func(*rpc.Response)

func SetResData(data proto.Message) ResponseOption {
	return func(p *rpc.Response) {
		p.Body, _ = proto.Marshal(data)
	}
}

func SetResErr(errmsg string) ResponseOption {
	return func(p *rpc.Response) {
		p.ErrMsg = errmsg
	}
}
