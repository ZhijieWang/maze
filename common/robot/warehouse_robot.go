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
	"fmt"
	"gonum.org/v1/gonum/graph"
	"maze/common/methods"
	"maze/common/trace"

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
func (r *simpleWarehouseRobot) Plan() {
	if r.act.GetType() == common.ActionTypeNull {
		if r.task == nil {
			if r.World.HasTasks() {
				t := r.World.GetNextTask()
				success, err := r.World.ClaimTask(t.GetTaskID(), r.id)

				if !success {
					panic(fmt.Sprintf("Failed to update task, %+v", err))
				} else {

				}
				r.act = methods.PlanTaskAction(r.World.GetGraph(), r.location, t)
				r.task = t
			}
		}
	}
}

func (r *simpleWarehouseRobot) Execute() common.Trace {
	var rTrace common.Trace
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
				rTrace = &trace.MoveTrace{
					RobotID:   r.ID(),
					Source:    r.location,
					Target:    n,
					Timestamp: r.tick,
				}
				r.location = n

			} else {
				rTrace = &trace.MoveTrace{
					RobotID:   r.ID(),
					Source:    r.location,
					Target:    n,
					Timestamp: r.tick,
				}
				r.location = n
				r.act = move
			}
		} else {
			move.SetStatus(common.EndStatus)
			rTrace = &trace.MoveTrace{
				RobotID:   r.ID(),
				Source:    r.location,
				Target:    r.location,
				Timestamp: r.tick,
			}
			r.act = r.act.GetChild()
		}

	case common.ActionTypeStartTask:
		r.act = r.act.GetChild()
		rTrace = trace.TaskExecutionTrace{}
	case common.ActionTypeEndTask:
		// mark task complete and remove self task\
		err := r.World.TaskUpdate(r.task.GetTaskID(), common.Completed)
		if err != nil {
			panic(err)
		}
		r.task = nil
		r.act = r.act.GetChild()
		rTrace = trace.TaskExecutionTrace{}
	case common.ActionTypeNull:
		// choose to remain on the same location, no move.
		rTrace = trace.TaskExecutionTrace{}
	default:
		// do nothing

	}
	return rTrace
}

// Run is a function that can be run in a concurrent way
func (r *simpleWarehouseRobot) Run() common.Trace {
	r.tick += 1
	r.Plan()
	return r.Execute()
}
func NewSimpleWarehouseRobot(id common.RobotID, location graph.Node, world common.World) *simpleWarehouseRobot {
	s := simpleWarehouseRobot{
		id,
		location,
		nil,
		nil,
		0,
		action.Null(),
		world,
	}
	return &s
}

func (r *simpleWarehouseRobot) Init() bool {
	return true
}
func (r *simpleWarehouseRobot) GetStatus() (common.Action, common.Task) {
	return r.act, r.task
}
