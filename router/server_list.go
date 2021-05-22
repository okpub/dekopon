package router

import (
	"math/rand"

	"github.com/okpub/dekopon/bean/message/club"
)

type ServerList struct {
	round int
	nodes []*club.Node
}

func (server *ServerList) AddNode(node *club.Node) {
	server.nodes = append(server.nodes, node)
}

//轮询
func (server *ServerList) GetRound() *club.Node {
	if server.round++; server.round >= len(server.nodes) {
		server.round = 0
	}
	return server.nodes[server.round]
}

//随机
func (server *ServerList) GetRandom() *club.Node {
	var index = rand.Int() % len(server.nodes)
	return server.nodes[index]
}

//hash获取
func (server *ServerList) GetHash(uid int) *club.Node {
	var index = uid % len(server.nodes)
	return server.nodes[index]
}
