package codec

import (
	"bytes"

	"github.com/skimmer/conn/packet"
)

type PomeloPacketDecoder struct{}

func NewPomeloPacketDecoder() PacketDecoder {
	return &PomeloPacketDecoder{}
}

func (*PomeloPacketDecoder) Decode(data []byte) (packets []*packet.Packet, err error) {
	var buf = bytes.NewBuffer(nil)
	buf.Write(data)

	for {
		if buf.Len() < HeadSize {
			break
		}
		var (
			header = buf.Next(HeadSize)
			n      = BytesToInt(header)
		)
		//big long
		if n > MaxPacketSize {
			err = ErrPacketSizeExcced
			break
		}
		if buf.Len() < n {
			break
		}
		var p = packet.New()
		if err = p.Unmarshal(buf.Next(n)); err == nil {
			packets = append(packets, p)
		} else {
			break
		}
	}
	return
}
