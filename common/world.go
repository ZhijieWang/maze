// Copyright © 2018 Zhijie (Bill) Wang <wangzhijiebill@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.Println("init started")
	simID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	outfile, _ := os.Create(simID.String() + ".log") // update path for your needs
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
func CreateWorld(numRobots int) World {
	w := &TransparentWorld{}
	var g *simple.WeightedUndirectedGraph
	g = simple.NewWeightedUndirectedGraph(1, 10000000)
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
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}

		w.robots = append(w.robots, &Robot{id: rID, location: g.Nodes()[r.Intn(len(g.Nodes()))]})
	}
	w.grid = g
	w.timestamp = 0
	return w
}

// World interface defines the behavior of World simulation
type World interface {
	Simulate(policy func(w World, robot *Robot, t int) Trace, graphUpdate func(world World, trace Trace), tGenerator func(maxT int, w World) []*Task)
	GetGraph() *simple.WeightedUndirectedGraph
	GetRobots() []*Robot
	EdgeWeightPropagation(start graph.Node, step, depth int)
	GetTasks() []*Task
	SetTasks(tasks []*Task)
}

// TransparentWorld is a data holder for simulation, robot have full visibility of the world and themselves
type TransparentWorld struct {
	timestamp int
	robots    []*Robot
	grid      *simple.WeightedUndirectedGraph
	Tasks     []*Task
}

//SetTasks add tasks to the queue
func (w *TransparentWorld) SetTasks(tasks []*Task) {
	w.Tasks = append(w.Tasks, tasks...)
}

//GetTasks return current task queue
func (w *TransparentWorld) GetTasks() []*Task {
	return w.Tasks
}

//GetGraph returns the underlying graph
func (w TransparentWorld) GetGraph() *simple.WeightedUndirectedGraph {
	return w.grid
}

// GetRobots returns a list of robot from underlying storage
func (w TransparentWorld) GetRobots() []*Robot {
	return w.robots
}

// Simulate is a step function for time synchronized simulation
func (w *TransparentWorld) Simulate(policy func(w World, robot *Robot, t int) Trace, graphUpdate func(world World, trace Trace), tGenerator func(maxT int, w World) []*Task) {
	w.timestamp++
	//	log.Printf("%d Tasks in current world \n", len(w.GetTasks()))
	w.SetTasks(tGenerator(50, w))

	for _, r := range w.robots {
		t := policy(w, r, w.timestamp)
		log.Printf("%+v\n", t)
		graphUpdate(w, t)
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
