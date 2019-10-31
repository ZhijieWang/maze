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
package methods

import (
	"errors"
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"math/rand"
	"maze/common"
	"maze/common/action"
	"maze/common/task"
	"maze/common/trace"
)

// TaskGenerator is the generator function for randomly producing tasks
func TaskGenerator(maxTasks int, w common.World) []common.Task {

	n := w.GetGraph().Nodes().Len()
	var tList []common.Task
	for i := 0; i < maxTasks; i++ {
		if rand.Intn(2) > 0 {
			uid, _ := uuid.NewUUID()
			tList = append(tList, task.TimePriorityTask{
				ID:          uid,
				Origin:      w.GetGraph().Node((int64)(rand.Intn(n))),
				Destination: w.GetGraph().Node((int64)(rand.Intn(n))),
			})
		}
	}
	return tList
}

func NoMove(r common.Robot, t int) common.Trace {
	return &trace.MoveTrace{
		RobotID:   r.ID(),
		Source:    r.Location(),
		Target:    r.Location(),
		Timestamp: t,
	}
}

// // GraphReWeightByRadiation is a graph weight propagation method to recalculate graph edge weight by radiation
// func GraphReWeightByRadiation(world World, trace Trace) {
// 	for _, i := range world.GetRobots() {
// 		world.EdgeWeightPropagation(i.location, 3, 1)
// 	}
// }

// RandMove is a basic function, robot takes a random move that it can move to.
// if there is only one path, robot will move
// this is stateless, regardless of previous move taken
func RandMove(w common.World, r common.Robot, t int) common.Trace {
	locs := w.GetGraph().From(r.Location().ID())

	bufs := graph.NodesOf(locs)

	rTrace := &trace.MoveTrace{
		RobotID:   r.ID(),
		Source:    r.Location(),
		Target:    bufs[rand.Intn(len(bufs))],
		Timestamp: t,
	}
	// r.Location() = trace.Target
	return rTrace

}
func NextMove(graph simple.DirectedGraph, start graph.Node) graph.Node {
	// path, ok =
	return graph.From(start.ID()).Node()
}
func GetPath(start, end common.Location, g graph.Graph) ([]graph.Node, error) {
	pt, ok := path.BellmanFordFrom(start, g)
	if ok {
		p, _ := pt.To(end.ID())
		if len(p) == 0 {
			panic("No Path")
		}
		return p[1:], nil
	} else {
		return nil, errors.New("negative cycle")
	}
}
func PlanTaskAction(g graph.Graph, location common.Location, task common.Task) common.Action {
	var start common.Action
	var current common.Action
	if location == task.GetOrigination() {
		start = action.CreateBeginTaskAction(location)
		current = start
	} else {
		start = action.CreateMoveAction(location, task.GetOrigination())

		pat, err := GetPath(location, task.GetOrigination(), g)
		if err != nil {
			panic(err)
		}
		start.(*action.MoveAction).Path = pat
		start.SetChild(action.CreateBeginTaskAction(location))
		current = start.GetChild()
	}

	p, _ := GetPath(task.GetOrigination(), task.GetDestination(), g)
	current.SetChild(action.CreateMoveActionWithPath(task.GetOrigination(), task.GetDestination(), p))
	current.GetChild().SetChild(action.CreateEndTaskAction(task.GetDestination()))
	current.GetChild().GetChild().SetChild(action.Null())

	return start
}

// SelectTaskByDistance returns a task from queue, and returns that task. If there is an error, return err
func SelectTaskByDistance(tm common.PassiveTaskManager, robot common.Robot, world common.World) (common.PriorityTask, []graph.Node, error) {
	tq := tm.GetAllTasks()
	if len(tq) == 0 {
		log.Println("No Tasks")
		return nil, nil, nil
	}
	//claiming task
	var tMin common.PriorityTask
	minWeight := 0.0
	var weight float64
	var p, pMin []graph.Node
	pt, _ := path.BellmanFordFrom(robot.Location(), world.GetGraph())
	for _, t := range tq {

		p, weight = pt.To(t.GetOrigination().ID())

		if weight < minWeight {
			minWeight = weight
			tMin = t.(task.TimePriorityTask)
			pMin = p
		}
	}
	// err := tm.ClaimTask(tMin, robot.ID())
	return tMin, pMin, nil
}
