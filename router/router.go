package router

import (
	"github.com/okpub/dekopon/actor"
	"github.com/okpub/dekopon/bean/message/club"
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

}

//通过路由规则获取服务器
func (r *Router) Broadcast(server, messageType int32, data interface{}) (err error) {

	return
}

//获取某集群
func (r *Router) GetCluster(messageType int32) (pids []actor.PID) {

	return
}

//通过hash获取服务器接口
func (r *Router) GetHash(uid int32, serverType int32) (pids actor.PID, err error) {

	return
}
