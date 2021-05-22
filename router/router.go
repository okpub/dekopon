package router

import (
	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/club"
	"github.com/okpub/dekopon/utils"
)

//路由规则
const (
	_ int = iota + 0
	ROUND
	RANDOM
	WEIGHT
	HASH
)

type Node struct {
	actor.PID
	*club.Node
}

type DNSRouter interface{}

//路由
type Router struct {
	nodes map[int32]*ServerList
}

func (r *Router) Add(node *club.Node) {
	if _, ok := r.nodes[node.MessageType]; !ok {
		r.nodes[node.MessageType] = &ServerList{}
	}
	r.nodes[node.MessageType].AddNode(node)
}

func (r *Router) Remove(id int32) {
	var node, ok = r.nodes[id]
	if ok {
		delete(r.nodes, id)
		node.Close()
	}
}

//通过路由规则获取服务器
func (r *Router) Broadcast(server, messageType int32, data interface{}) (err error) {
	if node, ok := r.nodes[server]; ok {
		err = node.Send(data)
	} else {
		err = utils.NilErr
		var nodes = r.GetCluster(messageType)
		for _, node := range nodes {
			err = node.Send(data)
		}
	}
	return
}

//获取某集群
func (r *Router) GetCluster(messageType int32) (pids []actor.PID) {
	for _, node := range r.nodes {
		if node.Node.MessageType == messageType {
			pids = append(pids, node.PID)
		}
	}
	return
}

//通过hash获取服务器接口
func (r *Router) GetHash(uid int32, serverType int32) (pids actor.PID, err error) {

	return
}
