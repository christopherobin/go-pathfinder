package pathfinder

type NodeSet map[int]Node

func NewNodeSet(size int) NodeSet {
	return make(NodeSet, size)
}

func (nodes NodeSet) RegisterNode(id int, connections []int, data interface{}) {
	nodes[id] = Node{id, connections, data}
}
