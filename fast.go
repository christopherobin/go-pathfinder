package pathfinder

type CheckFunc func(NodeSet, int) bool

// This function allows every systems to be visited in the FastRoute
func CheckAll(NodeSet, int) bool {
	return true
}

// This is the fast route solver, it count the amount of connections to the target while allowing the
// user to disable unwanted nodes using the checkFunc parameter, the difference compared to the
// slow solver is that you cannot prioritise systems and only accept/refuse a node, meaning
// that it has a higher chance of failing to solve a route when filtering nodes
func (nodes NodeSet) FastSolver(fromId, toId int, check CheckFunc) ([]int, error) {
	visited := map[int]int{}

	// Initialize the queue with the starting node
	queue := []Connection{{fromId, -1}}

	// Add the start node to the visited list with distance 0
	visited[fromId] = 0

	var current Connection

	// while there are elements in the queue
	for len(queue) > 0 {
		// Get first element and remove it from the queue
		current, queue = queue[0], queue[1:]

		// Add node to the visited list with distance +1 from parent
		if current.From != -1 {
			visited[current.To] = visited[current.From] + 1
		}

		// If we found our target, exit asap to start resolving the path
		if current.To == toId {
			break
		}

		// iterate on very connection
		for _, connectionId := range nodes[current.To].Connections {
			// If we didn't visit that node yet
			if _, visit := visited[connectionId]; visit == false {
				// Ask checkFunc if we allow that system
				if check(nodes, connectionId) {
					queue = append(queue, Connection{connectionId, current.To})
				}
				// Set the system to visited so that we won't add it twice (with a placeholder)
				visited[connectionId] = -1
			}
		}
	}

	// if no path was found, break here
	if current.To != toId {
		return nil, ErrNoPathFound
	}

	// now walk back from target to source
	path := []int{}
	walk := toId
	for {
		path = append([]int{walk}, path...)

		// we have arrived
		if walk == fromId {
			break
		}

		dist := visited[walk]
		for _, connId := range nodes[walk].Connections {
			// did we visit that guy in the original parsing?
			cdist, ok := visited[connId]
			if !ok {
				continue
			}
			// distance if current minus one, we are getting closer, walk from there next
			if cdist == (dist - 1) {
				walk = connId
			}
		}
	}

	return path, nil
}
