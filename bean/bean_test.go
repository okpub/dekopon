package bean

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/okpub/dekopon/bean/message/club"
	"github.com/okpub/dekopon/bean/message/login"
	"github.com/okpub/dekopon/bean/message/world"
)

func TestInit(t *testing.T) {
	fmt.Println(login.LoginReq{})
	fmt.Println(world.AddUserReq{})
	fmt.Println(club.RemoveNodeReq{})

	var req = &login.LoginReq{UserID: 101, AppID: 100, Pwd: "密码"}
	var body, _ = proto.Marshal(req)
	fmt.Println(body)

	var resp = &login.LoginReq{}
	proto.Unmarshal(body, resp)
	fmt.Println("啥书:", resp)
}
