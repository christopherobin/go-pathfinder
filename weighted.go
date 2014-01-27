package pathfinder

import (
	"math"
	"sort"
)

// Allows us to sort a list of nodes by weight
type ByWeight struct {
	nodes  []int
	weight map[int]float64
}

func (s ByWeight) Len() int      { return len(s.nodes) }
func (s ByWeight) Swap(i, j int) { s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i] }
func (s ByWeight) Less(i, j int) bool {
	var fromWeight, toWeight float64
	var present bool

	// If the first node doesn't have a weight, then it's position is always higher or equal
	// to the other node
	if fromWeight, present = s.weight[s.nodes[i]]; !present {
		return false
	}

	// If the second node doesn't have a weight, then we should move the first one to the
	// front
	if toWeight, present = s.weight[s.nodes[j]]; !present {
		return true
	}

	return fromWeight < toWeight
}

type WeightFunc func(nodes NodeSet, previous, current, to int) (float64, error)

// Give the same weight to every nodes, will return the same results as the FastSolver when using
// the checkAll function, if this function returns ErrInvalidNode then the system is removed from
// the solvable routes (can be used when solving a graph where a node is still being cached but
// considered as unpassable at the moment)
func WeightFuncConnections(nodes NodeSet, previous, current, to int) (float64, error) {
	return 1.0, nil
}

// This is the slow path solver, it uses a weight system instead of connection count and as such
// need to solve slightly more nodes before being sure of the result, it allows to solves
// some problems that the fast solver can't such as finding routes to a distant node where
// intermediate nodes are considered as very low priority.
//
// Be careful that your weight function doesn't return very large numbers, the solver will
// fail to solve correctly any route whose value is above math.Maxfloat64
func (nodes NodeSet) WeightedSolver(fromId, toId int, weightFunc WeightFunc) ([]int, error) {
	var previous int

	// This is a map that olds the weighted distance from the source for every system visited
	weights := map[int]float64{}
	// The maximum score is the highest number we can think of
	top := math.MaxFloat64
	queue := []Connection{}

	// Set from to be visited with distance 0
	weights[fromId] = 0.0

	// add the start node to the queue
	queue = append(queue, Connection{fromId, -1})
	var current Connection

	// while there are elements in the queue
	for len(queue) > 0 {
		// get first element and remote it from the queue
		current, queue = queue[0], queue[1:]

		// Iterate on very connection
		for _, connectionId := range nodes[current.To].Connections {
			localWeight, err := weightFunc(nodes, current.From, current.To, connectionId)

			// Ignore nodes that are set to invalid
			if err == ErrInvalidNode {
				continue
			}

			// Any other error we don't know, bail out
			if err != nil {
				return []int{}, err
			}

			// If we didn't visit that system yet
			newWeight := weights[current.To] + localWeight

			// Prune anything higher than the best route we found
			if newWeight > top {
				continue
			}

			// If we don't have a weight or if the new weight is lower than previously
			// found, try this new road
			if weight, ok := weights[connectionId]; !ok || newWeight < weight {
				weights[connectionId] = newWeight

				// If we are at our target, announce to everyone that we found a new
				// weight, or a weight lower than the previous one
				if connectionId == toId {
					top = newWeight
					// We can't break yet, there may be some shorter routes available
					continue
				}

				// Add it to the queue
				queue = append(queue, Connection{connectionId, current.To})
			}
		}
	}

	// now walk back from target to source
	path := []int{}
	walk := toId
	for {
		path = append([]int{walk}, path...)
		previous = walk

		// We have arrived back to where we are starting
		if walk == fromId {
			break
		}

		// Sort connections by weight ascending then iterate on them, the first valid connection
		// that we find that is valid we use for the next iteration
		sort.Sort(ByWeight{nodes[walk].Connections, weights})
		for _, connId := range nodes[walk].Connections {
			// did we visit that guy in the original parsing?
			_, ok := weights[connId]
			if !ok {
				continue
			}

			walk = connId
			break
		}

		// If the loop above didn't yield any result, then we failed to find a route
		if walk == previous {
			break
		}
	}

	// If the last position we checked is not the start system, then we failed to find a route
	if walk != fromId {
		return nil, ErrNoPathFound
	}

	// Return the route
	return path, nil
}
