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

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

//
//func init() {
//	log.SetPrefix("LOG: ")
//	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
//	log.Println("init started")
//	simID, err := uuid.NewUUID()
//	if err != nil {
//		log.Fatal(err)
//	}

//	outfile, _ := os.Create(simID.String() + ".log") // update path for your needs
//	log.SetOutput(outfile)
//
//}

//	1	- 	5	-	9
//	|	X	|		|
//	2	-	6		10
//  	|		|		|
//	3		7		11
//	|		|		|
//	4	- 	8 	-	12

func CreateBlankWorld() World {
	s := simpleWorld{}
	s.grid = simple.NewWeightedUndirectedGraph(1, 10000000)
	return &s
}

//CreateWorld generates a network of 12 nodes
func CreateWorld(numRobots int, tm TaskManager) World {
	w := simpleWorld{}
	var g *simple.WeightedUndirectedGraph = simple.NewWeightedUndirectedGraph(1, 10000000)
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
	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}

		w.AddRobot(NewSimpleRobot(rID,
			g.Nodes().Node()))
	}
	w.grid = g
	return &w
}

// World interface defines the behavior of World simulation
type World interface {
	GetGraph() graph.Graph
	GetRobots() []Robot
	UpdateRobot(Robot) bool
	GetTasks() []Task
	AddRobot(r Robot) bool
	SetTasks(tasks []Task) bool
	ClaimTask(tid TaskID, rid RobotID)
}

// simpleWorld is the base implementation of a fully visible world, backed with Gonum Simple Graph
type simpleWorld struct {
	robots []Robot
	tasks  []Task
	grid   *simple.WeightedUndirectedGraph
}

// SetTasks allows the new tasks to be added to the world
func (s *simpleWorld) SetTasks(tasks []Task) bool {
	s.tasks = append(s.tasks, tasks...)
	return true
}

// GetTasks allows the rerieval of tasks (available only)
func (s *simpleWorld) GetTasks() []Task {
	return s.tasks
}

// GetGraph allows the retrieval of world state. The current implementation returns the full world. This is where visibility can be implemented
func (s *simpleWorld) GetGraph() graph.Graph {
	return s.grid
}

// ClaimTask defines the mechanims that a Robot can claim a given task from the world
func (s *simpleWorld) ClaimTask(tid TaskID, rid RobotID) {
}

// GetRobots implemnts the fucntionality for retrieval of list of robots
func (s *simpleWorld) GetRobots() []Robot {
	return s.robots
}

// AddRobots add more robots to the stack
func (s *simpleWorld) AddRobots(robots []Robot) bool {
	s.robots = append(s.robots, robots...)
	return true
}
func (s *simpleWorld) AddRobot(robot Robot) bool {
	s.robots = append(s.robots, robot)
	return true
}

func (s *simpleWorld) UpdateRobot(that Robot) bool {
	for i, this := range s.robots {
		if this.ID() == that.ID() {
			s.robots[i] = that
			return true

		}
	}
	return false

}
