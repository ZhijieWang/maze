/*
 *  Copyright (c) 2019 Zhijie (Bill) Wang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package world

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"maze/common"
)

//	1	- 	5	-	9
//	|	X	|		|
//	2	-	6		10
//  	|		|		|
//	3		7		11
//	|		|		|
//	4	- 	8 	-	12

//CreateWorld generates a network of 12 nodes
func CreateWorld(tm common.TaskManager) common.World {
	w := simpleWorld{}
	var g = simple.NewWeightedUndirectedGraph(1, 10000000)
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
	w.grid = g
	w.tm = tm
	return &w
}

// simpleWorld is the base implementation of a fully visible world, backed with Gonum Simple Graph
type simpleWorld struct {
	robots []common.Robot
	tm     common.TaskManager
	grid   *simple.WeightedUndirectedGraph
}

func (s *simpleWorld) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {
	return s.tm.TaskUpdate(taskID, status)
}

// SetTasks allows the new tasks to be added to the world
func (s *simpleWorld) AddTasks(tasks []common.Task) bool {
	s.tm.AddTasks(tasks)
	return true
}

// GetAllTasks allows the retrieval of tasks (available only)
func (s *simpleWorld) GetAllTasks() []common.Task {
	return s.tm.GetAllTasks()
}

// GetAllTasks allows the retrieval of tasks (available only)
func (s *simpleWorld) GetTasks(n int) []common.Task {
	return s.tm.GetAllTasks()[:n]
}

// GetGraph allows the retrieval of world state. The current implementation returns the full world. This is where visibility can be implemented
func (s *simpleWorld) GetGraph() graph.Graph {
	return s.grid
}

// ClaimTask defines the mechanism that a Robot can claim a given task from the world
func (s *simpleWorld) ClaimTask(tid common.TaskID, rid common.RobotID) (success bool, err error) {
	err = s.tm.TaskUpdate(tid, common.Assigned)
	if err != nil {
		return false, err
	} else {
		return true, err
	}
}
func (s *simpleWorld) AddTask(t common.Task) bool {
	s.tm.AddTask(t)
	return true
}

// GetRobots implements the functionality for retrieval of list of robots
func (s *simpleWorld) GetRobots() []common.Robot {
	return s.robots
}

// AddRobots add more robots to the stack
func (s *simpleWorld) AddRobots(robots []common.Robot) bool {
	s.robots = append(s.robots, robots...)
	return true
}

func (s *simpleWorld) GetBroadcastInfo() interface{} {
	return struct{}{}
}
func (s *simpleWorld) GetNextTask() common.Task {
	return s.tm.GetNextTask()
}

// AddRobot function add individual robot to tracking on the world map
func (s *simpleWorld) AddRobot(robot common.Robot) bool {
	s.robots = append(s.robots, robot)
	return true
}

func (s *simpleWorld) UpdateRobot(that common.Robot) bool {
	for i, this := range s.robots {
		if this.ID() == that.ID() {
			s.robots[i] = that
			return true
		}
	}
	return false

}

func (s *simpleWorld) HasTasks() bool {
	return s.tm.HasTasks()
}
