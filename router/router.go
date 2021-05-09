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

type RouterFunc func(*club.Node) actor.PID

type DNSRouter interface {
}

//路由
type Router struct {
	RouterFunc
}

func (r *Router) Add(node *club.Node) {

}

func (r *Router) Remove(id int32) {

}

//通过路由规则获取服务器
func (r *Router) Exchange(routerType int, serverType int32) (pids actor.PID, err error) {

	return
}

//获取所有服务器
func (r *Router) GetAll() (pids []actor.PID) {

	return
}

//获取某集群
func (r *Router) GetCluster(serverType int32) (pids []actor.PID) {

	return
}

//通过服务器id获取
func (r *Router) Get(id int32) (pid actor.PID, err error) {

	return
}

//通过hash获取服务器接口
func (r *Router) GetHash(uid int32, serverType int32) (pids actor.PID, err error) {

	return
}
