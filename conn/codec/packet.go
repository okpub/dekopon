package codec

import (
	"errors"

	"github.com/okpub/dekopon/conn/packet"
)

// Codec constants.
const (
	HeadSize = 4 //单纯头长度 headSize
	//PayloadSize   = HeadSize + 4 //所有头长度 headSize+bodySize
	MaxPacketSize = 1 << 24 //最大包长度 16MB
)

// ErrPacketSizeExcced is the error used for encode/decode.
var ErrPacketSizeExcced = errors.New("codec: packet size exceed")

// 解包/读包
type PacketDecoder interface {
	Decode([]byte) ([]*packet.Packet, error)
}

// 打包
type PacketEncoder interface {
	Encode(interface{}) ([]byte, error)
}

func NewBytes(n int) []byte {
	return make([]byte, n)
}

func Die(err error) bool {
	return err != nil
}
