package message

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/okpub/dekopon/bean/message/common"
)

func TestInit(t *testing.T) {
	fmt.Println("message init")

	var m = &common.CustomMessage{
		Header: &common.MessageHeader{Cmd: 101},
	}
	var ss = &common.Session{Header: &common.SessionHeader{Unix: 101}}
	m.Body, _ = proto.Marshal(ss)

	var b1, _ = proto.Marshal(m)

	var m1 = &common.UserMessage{}
	var err = proto.Unmarshal(b1, m1)

	fmt.Println("用户消息", m1, err)
}
