// Copyright © 2019 Zhijie (Bill) Wang <wangzhijiebill@gmail.com>
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
package impl

import (
	"log"
	"maze/common"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// simpleWarehouseRobot is a data holder struct for robot
type simpleWarehouseRobot struct {
	// id is the UUID of the robot
	id common.RobotID
	// location represents the robot's current location on the graph
	location graph.Node
	// task represents the current work the robot is trying to carry out
	task common.Task
	// path is the current planned path to deliver the task
	path []graph.Node
}

// ID returns the robot UUID
func (r *simpleWarehouseRobot) ID() common.RobotID {
	return r.id
}

// Run is a function that can be run in a concurrent way
func (r *simpleWarehouseRobot) Run(w common.World, tm common.TaskManager) common.Trace {
	var tick int = 1
	if r.task == nil {
		if tm.HasTasks() {
			r.task = tm.GetTasks(1)[0]
			w.ClaimTask(r.task.GetTaskID(), r.ID())
			return common.Trace{
				RobotID:   r.ID(),
				Source:    r.location,
				Target:    r.task.GetDestination(),
				Timestamp: tick,
			}
		} else {
			return common.Trace{}
		}
	} else if r.location == r.task.GetDestination() {

		// at target location.
		// unset task from robot
		// update task to be done
		err := tm.TaskUpdate(r.task.GetTaskID(), common.Completed)
		if err != nil {
			return common.Trace{}
		}
		return common.Trace{}
	}
	// go to next location in path
	return common.Trace{}
	//r.localWorld = worldReader.Observe(r.location)
}
func NewSimpleWarehouseRobot(id common.RobotID, location graph.Node) common.Robot {
	s := simpleWarehouseRobot{
		id,
		location,
		nil,
		nil,
	}
	return &s
}

type WarehouseWorld struct {
	graph  *simple.DirectedGraph
	robots map[common.RobotID]common.Robot
}

func (w *WarehouseWorld) GetGraph() graph.Graph {
	return w.graph

}

func (w *WarehouseWorld) GetRobots() []common.Robot {
	values := make([]common.Robot, len(w.robots))
	for _, value := range w.robots {
		values = append(values, value)
	}
	return values
}
func (w *WarehouseWorld) GetTasks() []common.Task {
	panic("not implemented")
}

func (w *WarehouseWorld) SetTasks(tasks []common.Task) bool {
	panic("not implemented")
}

func (w *WarehouseWorld) ClaimTask(tid common.TaskID, rid common.RobotID) {
	panic("not implemented")
}

func CreateWarehouseWorld() *WarehouseWorld {

	w := &WarehouseWorld{
		simple.NewDirectedGraph(),
		make(map[common.RobotID]common.Robot),
	}
	for i := 1; i < 13; i++ {
		w.graph.AddNode(simple.Node(i))
	}
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(1), w.graph.Node(2)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(1), w.graph.Node(5)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(1), w.graph.Node(6)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(2), w.graph.Node(5)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(2), w.graph.Node(3)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(2), w.graph.Node(6)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(3), w.graph.Node(4)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(4), w.graph.Node(8)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(8), w.graph.Node(7)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(7), w.graph.Node(6)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(6), w.graph.Node(5)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(5), w.graph.Node(9)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(9), w.graph.Node(10)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(10), w.graph.Node(11)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(11), w.graph.Node(12)))
	w.graph.SetEdge(w.graph.NewEdge(w.graph.Node(12), w.graph.Node(8)))
	var numRobots int = 5

	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}

		w.AddRobot(NewSimpleWarehouseRobot(rID,
			w.graph.Nodes().Node()))
	}
	return w
}

func (w *WarehouseWorld) AddRobot(r common.Robot) bool {
	if _, ok := w.robots[r.ID()]; ok {
		// robot already in the track
		return false
	}
	return true
}
func (w *WarehouseWorld) UpdateRobot(r common.Robot) bool {

	if _, ok := w.robots[r.ID()]; ok {
		w.robots[r.ID()] = r
		return true
	}
	return false
}

type SimulatedTaskManager struct {
}

func CreateSimulatedTaskManager() *SimulatedTaskManager {
	return &SimulatedTaskManager{}
}
func (stm *SimulatedTaskManager) GetBroadcastInfo() interface{} {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) GetAllTasks() []common.Task {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) GetNext() common.Task {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) GetTasks(n int) []common.Task {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) AddTask(t common.Task) bool {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) AddTasks(tList []common.Task) bool {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) HasTasks() bool {
	panic("not implemented")
}
