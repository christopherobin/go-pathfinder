# go-pathfinder

A node based pathfinder in Go. I devised those solvers for fun, but they probably have a real
name, if you do know that name please open an issue so that I can rename/alias them.

## Install

`go get github.com/christopherobin/go-pathfinder`

## Example

```go
package main

import (
	"github.com/christopherobin/go-pathfinder"
	"log"
)

func main() {
	// create a small 6 nodes system, with some arbitrary data
	nodes := pathfinder.NewNodeSet(6)
	nodes.RegisterNode(1, []int{2, 3}, nil)
	nodes.RegisterNode(2, []int{1, 4}, nil)
	nodes.RegisterNode(3, []int{1, 4}, nil)
	nodes.RegisterNode(4, []int{2, 3, 5}, nil)
	nodes.RegisterNode(5, []int{4, 6}, nil)
	nodes.RegisterNode(6, []int{5}, nil)

	path, err := nodes.FastSolver(1, 6, pathfinder.CheckAll)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(path)
}
```

See `route_test.go` for more advanced usages on a larger data set.

## Benchmarks

Using the EVE online database (5428 nodes) and doing pathfinding from M-VACR to Shirshocin (~60 connections):

```
PASS
BenchmarkFast                1000       2404639 ns/op
BenchmarkWeighted            1000       2381592 ns/op
BenchmarkWeightedHighSec      500       5436040 ns/op
BenchmarkWeightedPrune       1000       2870525 ns/op
ok  	github.com/christopherobin/go-pathfinder	11.999s
```

* `BenchmarkWeightedHighSec`: Any system whose security rating is under 0.5 has a weight of 100 instead of 1
* `BenchmarkWeightedPrune`: 2 key systems are pruned from the path, forcing the solver to work his way around

The dataset used for the tests and benchmarks is in `test.json`

## Todo

* Name those solvers correctly
* Add heuristic support to the weighted solver