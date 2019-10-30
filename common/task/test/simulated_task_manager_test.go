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

package test_test

import (
	"github.com/google/uuid"
	"maze/common"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
	"testing"
)

var (
	w      common.World
	stm    *task.SimulatedTaskManager
	t1     *task.TimePriorityTask
	t2     *task.TimePriorityTask
	robots []common.Robot
)

func setup() {
	stm = task.CreateSimulatedTaskManager()
	w = world.CreateWarehouseWorldWithTaskManager(stm)

	robots = []common.Robot{}
	robots = append(robots, robot.NewSimpleWarehouseRobot(uuid.New(), w.GetGraph().Node(1), w))
	robots = append(robots, robot.NewSimpleWarehouseRobot(uuid.New(), w.GetGraph().Node(1), w))
}
func addT1() {
	t1 = task.NewTimePriorityTask()
	t1.Origin = w.GetGraph().Node(1)
	t1.Destination = w.GetGraph().Node(2)
	stm.AddTask(t1)
}
func addT2() {
	t2 = task.NewTimePriorityTask()
	t2.Origin = w.GetGraph().Node(1)
	t2.Destination = w.GetGraph().Node(6)
	stm.AddTask(t2)
}

func TestHasTasks(t *testing.T) {
	setup()
	addT1()
	addT2()

	// Expect HasTask to Be True
	if stm.HasTasks() && len(stm.GetAllTasks()) == 2 {
		// no problem
	} else {
		t.Errorf("The HasTasks should return true, given no claim has been done")
	}
	err := stm.TaskUpdate(t1.GetTaskID(), common.Assigned)
	err = stm.TaskUpdate(t1.GetTaskID(), common.Completed)
	if err != nil {
		t.Error("Update failed")
		t.FailNow()
	}
	// Expect HasTask to Be True
	if stm.HasTasks() && stm.ActiveCount() == 0 && stm.FinishedCount() == 1 {
		// no problem
	} else {
		t.Errorf("The HasTasks should return true, though 1 out of 2 is done")
	}
	err = stm.TaskUpdate(t2.GetTaskID(), common.Assigned)
	err = stm.TaskUpdate(t2.GetTaskID(), common.Completed)
	if err != nil {
		t.Error("Update failed")
		t.FailNow()
	}
	// Expect HasTask to Be True
	if stm.HasTasks() {
		t.Errorf("The HasTasks should return false, 2 out of 2 are done. %d active, %d finished, all remain %+v", stm.ActiveCount(), stm.FinishedCount(), stm.GetAllTasks())
	}
}

func TestRobotClaimTask(t *testing.T) {
	setup()
	addT1()
	id, _ := uuid.NewUUID()

	r := robot.NewSimpleWarehouseRobot(id, w.GetGraph().Node(1), w)
	_, rT := r.GetStatus()
	if rT != nil {
		t.Error("Shouldn't have any task on newly instantiated item")
	}
	r.Run()

	if stm.HasTasks() {
		t.Errorf("Robot failed to update task status")
		t.FailNow()
	}
	_, rT = r.GetStatus()
	if rT == nil {
		t.Error("Failed to claim task")
	}
}
