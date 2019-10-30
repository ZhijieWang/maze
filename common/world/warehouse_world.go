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

type WarehouseWorld struct {
	graph  *simple.UndirectedGraph
	robots map[common.RobotID]common.Robot
}

func (w *WarehouseWorld) GetGraph() graph.Graph {
	return w.graph

}

func (w *WarehouseWorld) GetRobots() []common.Robot {
	var values []common.Robot
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
		simple.NewUndirectedGraph(),
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
