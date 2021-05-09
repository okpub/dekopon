package codec

import (
	"fmt"

	"github.com/skimmer/conn/packet"
	"google.golang.org/protobuf/proto"
)

type PomeloPacketEncoder struct{}

func NewPomeloPacketEncoder() PacketEncoder {
	return &PomeloPacketEncoder{}
}

func (c *PomeloPacketEncoder) Encode(message interface{}) (data []byte, err error) {
	switch p := message.(type) {
	case proto.Message:
		if data, err = proto.Marshal(p); err == nil {
			data, err = c.Encode(packet.With(data))
		}
	case packet.Message:
		if data, err = p.Marshal(); err == nil {
			data, err = c.pack(data)
		}
	default:
		panic(fmt.Errorf("the message untype %T", message))
	}
	return
}

func (*PomeloPacketEncoder) pack(body []byte) (data []byte, err error) {
	if len(body) > MaxPacketSize {
		err = ErrPacketSizeExcced
	} else {
		data = make([]byte, HeadSize+len(body))
		//写入头
		//copy(data[:HeadSize], p.Header)
		//写入body长度
		IntToBytes(data[:HeadSize], len(body))
		//写入body
		copy(data[HeadSize:], body)
	}
	return
}
