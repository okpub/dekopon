package codec

import "encoding/binary"

func BytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

func IntToBytes(buf []byte, n int) {
	binary.BigEndian.PutUint32(buf, uint32(n))
}
