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
	"fmt"
	"log"
	"math/rand"
	"maze/common"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
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
	tick int
}

// ID returns the robot UUID
func (r *simpleWarehouseRobot) ID() common.RobotID {
	return r.id
}
func (r *simpleWarehouseRobot) Location() graph.Node {
	return r.location
}

//TaskMove is a movement policy for Task Oriented movement
func TaskMove(w common.World, r simpleWarehouseRobot, t int) common.Trace {
	log.Printf("Robot %s can see %d Tasks, current has %vi\n", r.id, len(w.GetTasks()), r.task)
	if r.task != nil {
		log.Printf("Robot %s is carrying out Task %+v\n", r.id, r.task)
		fmt.Printf("%+v\n", r.path)
		fmt.Printf("current location %s, task target location %s\n", r.Location(), r.task.GetDestination())
		trace := common.Trace{
			RobotID:   r.ID(),
			Source:    r.Location(),
			Target:    r.path[0],
			Timestamp: t,
		}
		r.location = r.path[0]
		if len(r.path) == 1 {
			log.Printf("Task %+v done by Robot %s\n", r.task, r.id)
			r.task = nil
			r.path = nil
		} else {
			r.path = r.path[1:]
		}
		return trace
	}
	tasks := w.GetTasks()
	if len(tasks) == 0 {
		log.Println("No Tasks")

		return common.NoMove(w, &r, t)
	}
	tMin := tasks[rand.Intn(len(tasks))]

	pt, _ := path.BellmanFordFrom(r.Location(), w.GetGraph())
	p, _ := pt.To(tMin.GetDestination().ID())
	r.path = p[1:]
	r.task = tMin
	return common.Trace{
		RobotID:   r.ID(),
		Source:    r.Location(),
		Target:    p[0],
		Timestamp: t,
	}
}

// Run is a function that can be run in a concurrent way
func (r *simpleWarehouseRobot) Run(w common.World, tm common.TaskManager) common.Trace {
	r.tick += 1
	// return TaskMove(w, *r, r.tick)
	if r.task == nil {
		if tm.HasTasks() {
			// there is at least a task to claim
			(r.task) = tm.GetNext()
			err := tm.TaskUpdate(r.task.GetTaskID(), common.Assigned)
			if err != nil {
				panic("Task Update Failed")
			}
			return common.Trace{
				RobotID:   r.ID(),
				Source:    r.location,
				Target:    r.task.GetDestination(),
				Timestamp: r.tick,
			}
		}
		// there is no task to do
	} else if r.location == r.task.GetDestination() {
		// at target location.
		// unset task from robot
		// update task to be done
		err := tm.TaskUpdate(r.task.GetTaskID(), common.Completed)

		if err != nil {
			panic("Failed to update task")
		}
		r.task = nil
		return common.Trace{
			RobotID:   r.id,
			Source:    r.location,
			Target:    graph.Empty.Node(),
			Timestamp: r.tick,
		}
	}
	// go to next location in path
	return common.Trace{
		RobotID:   r.id,
		Source:    r.location,
		Target:    graph.Empty.Node(),
		Timestamp: r.tick,
	}
	// r.localWorld = worldReader.Observe(r.location)
}
func NewSimpleWarehouseRobot(id common.RobotID, location graph.Node) common.Robot {
	s := simpleWarehouseRobot{
		id,
		location,
		nil,
		nil,
		0,
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
	values := []common.Robot{}
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

//	1	- 	5	-	9
//	|	X	|		|
//	2	-	6		10
//  	|		|		|
//	3		7		11
//	|		|		|
//	4	- 	8 	-	12
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
	w.robots[r.ID()] = r
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
	tasks   map[common.TaskID]common.Task
	active  map[common.TaskID]common.Task
	archive map[common.TaskID]common.Task
}

func CreateSimulatedTaskManager() *SimulatedTaskManager {
	return &SimulatedTaskManager{
		make(map[common.TaskID]common.Task),
		make(map[common.TaskID]common.Task),
		make(map[common.TaskID]common.Task),
	}
}
func (stm *SimulatedTaskManager) GetBroadcastInfo() interface{} {
	panic("not implemented")
}

func (stm *SimulatedTaskManager) GetAllTasks() []common.Task {
	values := []common.Task{}
	for _, t := range stm.tasks {
		values = append(values, t)
	}
	return values
}

func (stm *SimulatedTaskManager) GetNext() common.Task {
	if len(stm.tasks) == 0 {
		return nil
	}
	for _, v := range stm.tasks {
		if v.GetStatus() != common.Assigned {
			return v
		}

	}
	return nil
}

func (stm *SimulatedTaskManager) GetTasks(n int) []common.Task {
	values := make([]common.Task, n)
	for _, t := range stm.tasks {
		n -= 1
		values = append(values, t)
		if n == 0 {
			break
		}
	}
	return values
}

func (stm *SimulatedTaskManager) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {

	if status == common.Completed {
		stm.archive[taskID] = stm.tasks[taskID]
		delete(stm.tasks, taskID)
		return nil
	} else if status == common.Assigned {
		stm.active[taskID] = stm.tasks[taskID]
		delete(stm.tasks, taskID)
		return nil
	} else {
		return stm.tasks[taskID].UpdateStatus(status)
	}

}

func (stm *SimulatedTaskManager) AddTask(t common.Task) bool {
	if t.GetStatus() == common.Completed {
		return false
	} else {
		if _, ok := stm.tasks[t.GetTaskID()]; ok {
			// task already in the tracker
			// edge case, return false for now
			return false
		} else {
			stm.tasks[t.GetTaskID()] = t
			return true
		}
	}
}

func (stm *SimulatedTaskManager) AddTasks(tList []common.Task) bool {
	result := true
	for _, t := range tList {
		result = result && stm.AddTask(t)
	}
	return result
}

func (stm *SimulatedTaskManager) HasTasks() bool {
	return (0 != len(stm.tasks))
}
