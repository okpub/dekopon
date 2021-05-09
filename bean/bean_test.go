package bean

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/skimmer/bean/message/club"
	"github.com/skimmer/bean/message/login"
	"github.com/skimmer/bean/message/world"
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
