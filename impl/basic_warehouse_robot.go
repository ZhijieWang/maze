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
package impl

import (
	"errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"maze/common"
	"maze/common/action"
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
	act  common.Action

	common.World // a place to read world,
	common.TaskManager
}

// ID returns the robot UUID
func (r *simpleWarehouseRobot) ID() common.RobotID {
	return r.id
}
func (r *simpleWarehouseRobot) Location() graph.Node {
	return r.location
}

func (r *simpleWarehouseRobot) Stop() {

}
func (r *simpleWarehouseRobot) Plan(g graph.Graph) {
	r.act = PlanTaskAction(g, r.Location(), r.task)
}

func PlanTaskAction(g graph.Graph, location common.Location, task common.Task) common.Action {
	var start common.Action
	var current common.Action
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

func (r *simpleWarehouseRobot) Execute(g graph.Graph, tm common.TaskManager) (graph.Node, common.Action) {

	switch r.act.GetType() {
	case common.ActionTypeMove:
		move := r.act.(*action.MoveAction)
		move.SetStatus(common.ActiveStatus)
		if len(move.Path) > 0 {
			n := move.Path[0]

			move.Path = move.Path[1:]
			if len(move.Path) == 0 {
				move.SetStatus(common.EndStatus)
				r.act = move.GetChild()
				r.location = n
			} else {
				r.location = n
				r.act = move
			}
		} else {
			move.SetStatus(common.EndStatus)
			r.act = r.act.GetChild()
		}
	case common.ActionTypeStartTask:
		tm.TaskUpdate(r.task.GetTaskID(), common.Assigned)
		r.act = r.act.GetChild()
	case common.ActionTypeEndTask:
		// mark task complete and remove self task\

		tm.TaskUpdate(r.task.GetTaskID(), common.Completed)
		r.task = nil
		r.act = r.act.GetChild()
	case common.ActionTypeNull:
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
func (r *simpleWarehouseRobot) Run() common.Trace {
	source := r.location
	r.tick += 1
	if r.act.GetType() == common.ActionTypeNull {
		if r.task == nil {
			if r.TaskManager.HasTasks() {
				t := r.TaskManager.GetNext()
				err := r.TaskManager.TaskUpdate(t.GetTaskID(), common.Assigned)
				if err != nil {
					panic("Failed to update task")
				}
				r.act = PlanTaskAction(r.World.GetGraph(), r.location, t)
				r.task = t
			}
		}
	}

	n, act := r.Execute(r.World.GetGraph(), r.TaskManager)
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
func NewSimpleWarehouseRobot(id common.RobotID, location graph.Node, world common.World, manager common.TaskManager) *simpleWarehouseRobot {
	s := simpleWarehouseRobot{
		id,
		location,
		nil,
		nil,
		0,
		action.Null(),
		world,
		manager,
	}
	return &s
}

func (r *simpleWarehouseRobot) Init() bool {
	return true
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

	switch status {

	case common.Completed:
		if t, ok := stm.active[taskID]; ok {
			stm.archive[taskID] = t
			delete(stm.active, taskID)
			return nil
		} else {
			return errors.New("status can't jump from UnAssigned to Completed")
		}

	case common.Assigned:
		if t, ok := stm.tasks[taskID]; ok {
			stm.active[taskID] = t
			delete(stm.tasks, taskID)
			return nil
		} else {
			return errors.New("task not found")
		}

	default:
		return nil
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
