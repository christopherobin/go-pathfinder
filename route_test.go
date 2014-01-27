package pathfinder

import (
	"encoding/json"
	"os"
	"testing"
)

// For testing we use the EVE Online map, it contains several thousand nodes
// and allow us to test several interesting weight functions
type EVESystem struct {
	Id          int     `json:"id"`          // The internal ID used by CCP for this system
	Name        string  `json:"name"`        // The name of the system
	Security    float64 `json:"security"`    // The security rating, ranges from -1.0 to 1.0
	Connections []int   `json:"connections"` // The list of connections, could be found using the stargates but faster that way
}

// A few variables that point to known systems
var TestSystems = map[string]int{
	"Amarr":      30002187,
	"D85-VD":     30003345,
	"IMK-K1":     30000895,
	"M-VACR":     30004213,
	"Shirshocin": 30004233,
	"Youl":       30003493,
}

var testFile = "test.json"

func SetupNodes() (NodeSet, error) {
	file, err := os.Open(testFile)
	if err != nil {
		return nil, err
	}

	var systems map[string]EVESystem

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&systems)

	if err != nil {
		return nil, err
	}

	nodes := NewNodeSet(len(systems))

	for _, system := range systems {
		nodes.RegisterNode(system.Id, system.Connections, interface{}(system))
	}

	return nodes, nil
}

func TestFastSolver(t *testing.T) {
	nodes, err := SetupNodes()
	if err != nil {
		t.Fatal(err)
	}

	path, err := nodes.FastSolver(TestSystems["Amarr"], TestSystems["Youl"], CheckAll)
	if err != nil {
		t.Error(err)
	}

	if len(path) != 5 {
		t.Errorf("Path should have 5 jumps, %d jumps returned", len(path))
	}

	path2, err := nodes.FastSolver(TestSystems["M-VACR"], TestSystems["Shirshocin"], CheckAll)
	if err != nil {
		t.Error(err)
	}

	if len(path2) != 60 {
		t.Errorf("Path should have 60 jumps, %d jumps returned", len(path2))
	}
}

func TestWeightedSolver(t *testing.T) {
	nodes, err := SetupNodes()
	if err != nil {
		t.Fatal(err)
	}

	path, err := nodes.WeightedSolver(TestSystems["Amarr"], TestSystems["Youl"], WeightFuncConnections)
	if err != nil {
		t.Error(err)
	}

	if len(path) != 5 {
		t.Errorf("Path should have 5 jumps, %d jumps returned", len(path))
	}

	pruneFunc := func(nodes NodeSet, previous, current, to int) (float64, error) {
		// We take the second and fourth system found in the original resolution and prune them
		if to == path[1] || to == path[3] {
			return 0.0, ErrInvalidNode
		}
		// Once pruned, we use the basic weight function
		return WeightFuncConnections(nodes, previous, current, to)
	}
	prunedPath, err := nodes.WeightedSolver(TestSystems["Amarr"], TestSystems["Youl"], pruneFunc)

	if len(prunedPath) != 5 {
		t.Errorf("Path should have 5 jumps, %d jumps returned", len(prunedPath))
	}

	// check that the system we removed from the resolution were really removed
	for _, id := range prunedPath {
		if id == path[1] || id == path[3] {
			t.Errorf("A pruned system was found in returned path: %d", id)
		}
	}
}

// Benchmarks
func BenchmarkFast(b *testing.B) {
	nodes, err := SetupNodes()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// fast route solver
		_, err := nodes.FastSolver(TestSystems["M-VACR"], TestSystems["Shirshocin"], CheckAll)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWeighted(b *testing.B) {
	nodes, err := SetupNodes()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = nodes.WeightedSolver(TestSystems["M-VACR"], TestSystems["Shirshocin"], WeightFuncConnections)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWeightedHighSec(b *testing.B) {
	nodes, err := SetupNodes()
	if err != nil {
		b.Fatal(err)
	}

	securityWeightFunc := func(nodes NodeSet, from, current, to int) (float64, error) {
		// This callback gives a high weight to any system having a security inferior to 0.5 (in eve, security ranges
		// from -1.0 to 1.0, anything under 0.5 is lawless) but doesn't prune them, using them if it is the only
		// available path (or if the secure path takes more than 100 jumps per low-sec sector to get around)
		toSystem := nodes[to].Data.(EVESystem)
		if toSystem.Security < 0.5 {
			return 100.0, nil
		}

		return 1.0, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = nodes.WeightedSolver(TestSystems["M-VACR"], TestSystems["Shirshocin"], securityWeightFunc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWeightedPrune(b *testing.B) {
	nodes, err := SetupNodes()
	if err != nil {
		b.Fatal(err)
	}

	pruneFunc := func(nodes NodeSet, previous, current, to int) (float64, error) {
		// With this one we prune 2 systems that normaly would be used for our test path,
		// forcing the system to work around those
		if to == TestSystems["D85-VD"] || to == TestSystems["IMK-K1"] {
			return 0.0, ErrInvalidNode
		}
		// Once pruned, we use the basic weight function
		return WeightFuncConnections(nodes, previous, current, to)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = nodes.WeightedSolver(TestSystems["M-VACR"], TestSystems["Shirshocin"], pruneFunc)
		if err != nil {
			b.Fatal(err)
		}
	}
}
