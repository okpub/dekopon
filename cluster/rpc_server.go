package cluster

import (
	"github.com/skimmer/bean/message/common"
)

//应答服务器
type RPCServer struct {
}

func (rcp *RPCServer) Request(req *common.Request) (res *common.Response, err error) {

	return
}

func (rcp *RPCServer) Start() (err error) {

	return
}

func (rcp *RPCServer) Disconnect() {

}
