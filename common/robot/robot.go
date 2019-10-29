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
package robot

import (
	"gonum.org/v1/gonum/graph"
	"maze/common"
	"maze/common/action"
)

// SimpleRobot is a data holder struct for robot
type simpleRobot struct {
	// id is the UUID of the robot
	id common.RobotID
	// location represents the robot's current location on the graph
	location graph.Node
	// task represents the current work the robot is trying to carry out
	task common.Task
	// path is the current planned path to deliver the task
	path []graph.Node

	common.World
	common.TaskManager
}

// ID returns the robot UUID
func (r *simpleRobot) ID() common.RobotID {
	return r.id
}

// Run is a function that can be run in a concurrent way
func (r *simpleRobot) Run() common.Trace {

	var tick int = 1
	if r.task == nil {
		if r.TaskManager.HasTasks() {
			r.task = r.TaskManager.GetTasks(1)[0]

			r.World.ClaimTask(r.task.GetTaskID(), r.ID())
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
		return common.Trace{}
	}
	// go to next location in path
	return common.Trace{}
	//r.localWorld = worldReader.Observe(r.location)

}
func (r *simpleRobot) Location() graph.Node {
	return r.location
}
func NewSimpleRobot(id common.RobotID, location graph.Node, world common.World, manager common.TaskManager) common.Robot {
	s := simpleRobot{
		id,
		location,
		nil,
		nil,
		world,
		manager,
	}
	return &s
}

func (r *simpleRobot) Plan(g graph.Graph) {

}
func (r *simpleRobot) Init() bool {
	return true
}
func (r *simpleRobot) Execute(g graph.Graph, tm common.TaskManager) (graph.Node, common.Action) {
	return r.Location(), action.Null()
}
