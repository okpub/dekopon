package actor

type Node struct {
	parent PID
	self   PID
	system ActorSystem
}

func (node *Node) Self() PID           { return node.self }
func (node *Node) Parent() PID         { return node.parent }
func (node *Node) System() ActorSystem { return node.system }

func (node *Node) ChildOf(pid PID) infoPart {
	return &Node{parent: node.self, self: pid, system: node.system}
}
