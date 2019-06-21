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
	"os"

	"github.com/google/uuid"
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
	w := NewSimpleWorld()
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
	//r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < numRobots; i++ {
		//rID, err := uuid.NewUUID()
		//	if err != nil {
		//		log.Fatal(err)
		//	}

		//		w.robots = append(w.robots, &Robot{id: rID, location: g.Nodes()[r.Intn(len(g.Nodes()))]})
	}
	w.grid = g
	//w.timestamp = 0
	return &w
}

// World interface defines the behavior of World simulation
type World interface {
	GetGraph() *simple.WeightedUndirectedGraph
	GetRobots() []Robot
	//EdgeWeightPropagation(start graph.Node, step, depth int)
	GetTasks() []Task
	SetTasks(tasks []Task) bool
	ClaimTask(tid TaskID, rid RobotID)
}

// SimpleWorld is the base implementation of a fully visible world, backed with Gonum Simple Graph
type SimpleWorld struct {
	robots []Robot
	tasks  []Task
	grid   *simple.WeightedUndirectedGraph
}

// SetTasks allows the new tasks to be added to the world
func (s *SimpleWorld) SetTasks(tasks []Task) bool {
	s.tasks = append(s.tasks, tasks...)
	return true
}

// GetTasks allows the rerieval of tasks (available only)
func (s *SimpleWorld) GetTasks() []Task {
	return s.tasks
}

// GetGraph allows the retrieval of world state. The current implementation returns the full world. This is where visibility can be implemented
func (s *SimpleWorld) GetGraph() *simple.WeightedUndirectedGraph {
	return s.grid
}

// NewSimpleWorld is the constructor for simple world case
func NewSimpleWorld() SimpleWorld {
	return SimpleWorld{}
}

// ClaimTask defines the mechanims that a Robot can claim a given task from the world
func (s *SimpleWorld) ClaimTask(tid TaskID, rid RobotID) {
}

// GetRobots implemnts the fucntionality for retrieval of list of robots
func (s *SimpleWorld) GetRobots() []Robot {
	return s.robots
}

// AddRobots add more robots to the stack
func (s *SimpleWorld) AddRobots(robots []Robot) bool {
	s.robots = append(s.robots, robots...)
	return true
}
