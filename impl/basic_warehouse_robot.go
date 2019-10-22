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
	"errors"
	"log"
	"maze/common"
	"maze/common/action"

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
	act  action.Action
}

// ID returns the robot UUID
func (r *simpleWarehouseRobot) ID() common.RobotID {
	return r.id
}
func (r *simpleWarehouseRobot) Location() graph.Node {
	return r.location
}

//TaskMove is a movement policy for Task Oriented movement
// func TaskMove(w common.World, tm common.TaskManager, r simpleWarehouseRobot, t int) common.Trace {
// 	if r.task != nil {
// 		trace := common.Trace{
// 			RobotID:   r.ID(),
// 			Source:    r.Location(),
// 			Target:    r.path[0],
// 			Timestamp: t,
// 		}
// 		r.location = r.path[0]
// 		if len(r.path) == 1 {
// 			r.task = nil
// 			r.path = nil
// 		} else {
// 			r.path = r.path[1:]
// 		}
// 		return trace
// 	}
// 	tasks := w.GetTasks()
// 	if len(tasks) == 0 {
// 		log.Println("No Tasks")

// 		return common.NoMove(w, &r, t)
// 	}
// 	tMin := tasks[rand.Intn(len(tasks))]

// pt, _ := path.BellmanFordFrom(r.Location(), w.GetGraph())
// 	p, _ := pt.To(tMin.GetDestination().ID())
// 	r.path = p[1:]
// 	r.task = tMin
// 	return common.Trace{
// 		RobotID:   r.ID(),
// 		Source:    r.Location(),
// 		Target:    p[0],
// 		Timestamp: t,
// 	}
// }

func PlanTaskAction(g graph.Graph, location common.Location, task common.Task) action.Action {
	var start action.Action
	var current action.Action
	if location == task.GetOrigination() {
		start = action.CreateBeginTaskAction(location)
		current = start
	} else {
		start = action.CreateMoveAction(location, task.GetOrigination())
		start.(*action.MoveAction).Path, _ = GetPath(location, task.GetOrigination(), g)
		start.SetChild(action.CreateBeginTaskAction(location))
		current = start.GetChild()
	}
	p, _ := GetPath(task.GetOrigination(), task.GetDestination(), g)
	current.SetChild(action.CreateMoveActionWithPath(task.GetOrigination(), task.GetDestination(), p))
	current.GetChild().SetChild(action.CreateEndTaskAction(task.GetDestination()))
	current.GetChild().GetChild().SetChild(action.Null())

	return start
}

func NextMove(graph simple.DirectedGraph, start graph.Node) graph.Node {
	// path, ok =
	return graph.From(start.ID()).Node()
}

func (r *simpleWarehouseRobot) Execute(g graph.Graph, tm common.TaskManager) (graph.Node, action.Action) {

	switch r.act.GetType() {
	case action.ActionTypeMove:
		move := r.act.(*action.MoveAction)
		move.SetStatus(action.ActiveStatus)
		if len(move.Path) > 0 {
			n := move.Path[0]

			move.Path = move.Path[1:]
			if len(move.Path) == 0 {
				move.SetStatus(action.EndStatus)
				r.act = move.GetChild()
				r.location = n
			} else {
				r.location = n
				r.act = move
			}
		} else {
			move.SetStatus(action.EndStatus)
			r.act = r.act.GetChild()
		}
	case action.ActionTypeStartTask:
		tm.TaskUpdate(r.task.GetTaskID(), common.Assigned)
		r.act = r.act.GetChild()
	case action.ActionTypeEndTask:
		// mark task complete and remove self task
		tm.TaskUpdate(r.task.GetTaskID(), common.Completed)
		r.task = nil
		r.act = r.act.GetChild()
	case action.ActionTypeNull:
		// choose to remain on the same location, no move.
		return r.Location(), r.act
	default:
		// do nothing

	}
	return r.location, r.act
}
func GetPath(start, end common.Location, g graph.Graph) ([]graph.Node, error) {
	pt, ok := path.BellmanFordFrom(start, g)
	if ok {
		p, _ := pt.To(end.ID())

		return p[1:], nil
	} else {
		return nil, errors.New("no positive cycle")
	}
}

// Run is a function that can be run in a concurrent way
func (r *simpleWarehouseRobot) Run(w common.World, tm common.TaskManager) common.Trace {
	source := r.location
	r.tick += 1
	if r.act.GetType() == action.ActionTypeNull {
		if r.task == nil {
			if tm.HasTasks() {
				t := tm.GetNext()
				err := tm.TaskUpdate(t.GetTaskID(), common.Assigned)
				if err != nil {
					panic("Failed to update task")
				}
				r.act = PlanTaskAction(w.GetGraph(), r.location, t)
				r.task = t
			}
		}
	}

	n, act := r.Execute(w.GetGraph(), tm)
	r.act = act
	trace := common.Trace{
		RobotID:   r.ID(),
		Source:    source,
		Target:    r.Location(),
		Timestamp: r.tick,
	}
	r.location = n
	return trace
	// r.localWorld = worldReader.Observe(r.location)
}
func NewSimpleWarehouseRobot(id common.RobotID, location graph.Node) common.Robot {
	s := simpleWarehouseRobot{
		id,
		location,
		nil,
		nil,
		0,
		action.Null(),
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
			w.graph.Node(1)))
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
	return len(stm.tasks) != 0
}
func (stm *SimulatedTaskManager) FinishedCount() int {
	return len(stm.archive)
}

func (stm *SimulatedTaskManager) ActiveCount() int {
	return len(stm.active)
}