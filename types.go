package pathfinder

import "errors"

// A node is an entity of the graph, the Id must be uniques among
// every nodes in the graph or bad things will happen
type Node struct {
	Id          int         // We need a unique ID for each node
	Connections []int       // The list of nodes this node can access
	Data        interface{} // Custom interface that can be linked to a node
}

// 2 small helper types for code lisibility
type Nodes map[int]Node
type Connection struct {
	To   int
	From int
}

var ErrInvalidNode = errors.New("Invalid node")
var ErrNoPathFound = errors.New("No valid path found")
