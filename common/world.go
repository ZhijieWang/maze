package common

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"os"
)
func init(){
    log.SetPrefix("LOG: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
    log.Println("init started")
simID, err:=uuid.NewUUID()
	if (err!=nil){
		log.Fatal(err)
	}
    outfile, _ := os.Create(simID.String()+".log") // update path for your needs
    log.SetOutput(outfile)

}
//	1	- 	5	-	9
//	|	X	|		|
//	2	-	6		10
//  	|		|		|
//	3		7		11
//	|		|		|
//	4	- 	8 	-	12

//CreateWorld generates a network of 12 nodes
func CreateWorld(numRobots int, concurrent bool) World {
	w := &TransparentWorld{}
	w.Concurrency = concurrent

	g := simple.NewWeightedUndirectedGraph(1, 10000000)
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
		rID, err:= uuid.NewUUID()
	if (err!= nil){
		log.Fatal(err)
	}

		w.robots = append(w.robots, Robot{id: rID,location: g.Nodes()[r.Intn(len(g.Nodes()))]})
	}
	w.grid = g
	w.timestamp = 0
	return w
}

// Robot is a data holder struct for robot
type Robot struct {
	id uuid.UUID
	location graph.Node
}
// ID returns the robot UUID
func (r *Robot)ID() uuid.UUID{
	return r.id
}

// World interface defines the behavior of World simulation
type World interface {
	Simulate(policy func(w World, robot Robot, t int) Trace, graphUpdate func(world World, trace Trace))
	GetGraph() *simple.WeightedUndirectedGraph
	GetRobots() []Robot
	EdgeWeightPropagation(start graph.Node, step, depth int)
}

// TransparentWorld is a data holder for simulation, robot have full visibility of the world and themselves
type TransparentWorld struct {
	timestamp   int
	robots      []Robot
	grid        *simple.WeightedUndirectedGraph
	Concurrency bool
}

// Trace is data structure to hold data that can be used for path planning
type Trace struct {
	RobotID  uuid.UUID
	Source    graph.Node
	Target    graph.Node
	Timestamp int
}

//GetGraph returns the underlying graph
func (w TransparentWorld) GetGraph() *simple.WeightedUndirectedGraph {
	return w.grid
}

// RandMove is a basic function, robot takes a random move that it can move to.
// if there is onlyone path, robot will move
// this is stateless, regardless of previous move taken
func RandMove(w World, r Robot, t int) Trace {
	locs := w.GetGraph().From(r.location.ID())
	trace := Trace{
		RobotID:	r.ID(),
		Source:    r.location,
		Target:    locs[rand.Intn(len(locs))],
		Timestamp: t,
	}
	r.location = trace.Target
	return trace
}

// GetRobots returns a list of robot from underlying storage
func (w TransparentWorld) GetRobots() []Robot {
	return w.robots
}

// Simulate is a step function for time synchronized simulation
func (w TransparentWorld) Simulate(policy func(w World, robot Robot, t int) Trace, graphUpdate func(world World, trace Trace)) {
	w.timestamp++
	for _, r := range w.robots {
		t := policy(w, r, w.timestamp)
		log.Printf("%+v\n", t)
		graphUpdate(&w, t)
	}
}

// EdgeWeightPropagation is the edge weight update function
func (w TransparentWorld) EdgeWeightPropagation(start graph.Node, steps, depth int) {
	if steps > depth {
		nodes := w.grid.From(start.ID())
		for _, n := range nodes {
			e := w.grid.WeightedEdgeBetween(start.ID(), n.ID())
			w.grid.SetWeightedEdge(w.grid.NewWeightedEdge(e.From(), e.To(), e.Weight()-float64(1.0/float64(depth*depth))))
			w.EdgeWeightPropagation(n, steps, depth+1)

		}

	}
}

// GraphReWeightByRadiation is a graph weight propagation method to recalculate graph edge weight by radiation
func GraphReWeightByRadiation(world World, trace Trace) {
	for _, i := range world.GetRobots() {
		world.EdgeWeightPropagation(i.location, 3, 1)
	}
}

// Print is a printing utilito
func (w TransparentWorld) Print() {
	fmt.Println(w.grid)
}
