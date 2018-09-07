package common

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

//	1	- 	5	-	9
//	|	X	|		|
//	2	-	6		10
//  	|		|		|
//	3		7		11
//	|		|		|
//	4	- 	8 	-	12

//CreateWorld generates a network of 12 nodes
func CreateWorld(numRobots int, concurrent bool) *World {
	w := &World{}
	w.Concurrency = concurrent

	g := NewConcurrentWeightedUndirectedGraph(1, 10000000)
	for i := 1; i < 13; i++ {
		g.AddNode(simple.Node(i))
	}
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(1), simple.Node(2), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(1), simple.Node(5), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(1), simple.Node(6), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(2), simple.Node(5), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(2), simple.Node(3), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(2), simple.Node(6), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(3), simple.Node(4), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(4), simple.Node(8), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(8), simple.Node(7), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(7), simple.Node(6), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(6), simple.Node(5), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(5), simple.Node(9), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(9), simple.Node(10), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(10), simple.Node(11), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(11), simple.Node(12), 1))
	g.SetWeightedEdge(g.NewWeightedEdge(simple.Node(12), simple.Node(8), 1))
	//randomly assign x robots to positions
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < numRobots; i++ {
		w.robots = append(w.robots, Robot{location: g.Nodes()[r.Intn(len(g.Nodes()))]})
	}
	w.grid = g
	w.timestamp = 0
	return w
}

// Robot is a data holder struct for robot
type Robot struct {
	location graph.Node
}

// World is a data holder for simulation
type World struct {
	timestamp   int
	robots      []Robot
	grid        *ConcurrentWeightedUndirectedGraph
	Concurrency bool
}

// Trace is data structure to hold data that can be used for path planning
type Trace struct {
	Source    graph.Node
	Target    graph.Node
	Timestamp int
}

//RandMove is a basic function, robot takes a random move that it can move to.
// if there is onlyone path, robot will move
// this is stateless, regardless of previous move taken
func RandMove(w World, r Robot, t int) Trace {
	locs := w.grid.From(r.location.ID())
	trace := Trace{
		Source:    r.location,
		Target:    locs[rand.Intn(len(locs))],
		Timestamp: t,
	}
	r.location = trace.Target
	return trace
}

// Simulate is a step function for time synchronized simulation
func (w World) Simulate(policy func(w World, robot Robot, t int) Trace, graphUpdate func(world *World, trace Trace)) {
	w.timestamp++
	for _, r := range w.robots {

		t := policy(w, r, w.timestamp)
		graphUpdate(&w, t)

	}
	//	for _, edge := range w.grid.WeightedEdges() {
	//		fmt.Printf("%s %s %f\n", edge.From(), edge.To(), edge.Weight())
	//	}

}

// EdgeWeightPropagation is the edge weight update function
func (w World) EdgeWeightPropagation(start graph.Node, steps, depth int) {
	if steps > depth {
		nodes := w.grid.From(start.ID())
		for _, n := range nodes {
			go func() {
				w.UpdateWeight(w.grid.WeightedEdgeBetween(start.ID(), n.ID()), float64(1.0/float64(depth*depth)))
				w.EdgeWeightPropagation(n, steps, depth+1)
			}()
		}

	}
}

// GraphReWeightByRadiation is a graph weight propagation method to recalculate graph edge weight by radiation
func GraphReWeightByRadiation(world *World, trace Trace) {
	for _, i := range world.robots {
		world.EdgeWeightPropagation(i.location, 3, 1)
	}
}

// UpdateWeight is a short hand for update edge weight
func (w World) UpdateWeight(e graph.WeightedEdge, weightDelta float64) {

	w.grid.SetWeightedEdge(w.grid.NewWeightedEdge(e.From(), e.To(), e.Weight()-weightDelta))

}

// TestUpdate is a debug intermediar to make sure access pointers are defined correctly
func (w World) TestUpdate(x int64, y int64) {
	w.grid.SetWeightedEdge(w.grid.NewWeightedEdge(w.grid.Node(x), w.grid.Node(y), 100))
	fmt.Println(w.grid)
}

// Print is a printing utility
func (w World) Print() {
	fmt.Println(w.grid)
}
