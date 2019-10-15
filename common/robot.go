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
package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// RobotID is an alias to UUID for disambiguition purpose
type RobotID = uuid.UUID
type Robot interface {
	ID() RobotID
	Run(w World, tm TaskManager) Trace
}

// SimpleRobot is a data holder struct for robot
type simpleRobot struct {
	// id is the UUID of the robot
	id RobotID
	// location represents the robot's current location on the graph
	location graph.Node
	// task represents the current work the robot is trying to carry out
	task Task
	// path is the current planned path to deliver the task
	path []graph.Node
}

// ID returns the robot UUID
func (r *simpleRobot) ID() RobotID {
	return r.id
}

// Run is a function that can be run in a concurrent way
func (r *simpleRobot) Run(w World, tm TaskManager) Trace {
	var tick int = 1
	if r.task == nil {
		if tm.HasTasks() {
			r.task = tm.GetTasks(1)[0]
			w.ClaimTask(r.task.GetTaskID(), r.ID())
			return Trace{
				RobotID:   r.ID(),
				Source:    r.location,
				Target:    r.task.GetDestination(),
				Timestamp: tick,
			}
		} else {
			return Trace{}
		}
	} else if r.location == r.task.GetDestination() {
		// at target location.
		// unset task from robothttps://github.com/int3h/SublimeFixMacPath
		// update task to be done

		//		err := w.TaskUpdate(r.task.GetTaskID(), Completed)
		//		if err != nil {
		//			return Trace{}
		//		}
		return Trace{}
	}
	// go to next location in path
	return Trace{}
	//r.localWorld = worldReader.Observe(r.location)

}
func NewSimpleRobot(id RobotID, location graph.Node) Robot {
	s := simpleRobot{
		id,
		location,
		nil,
		nil,
	}
	return &s
}
