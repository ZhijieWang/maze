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
	"maze/common"
	"maze/common/action"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"gonum.org/v1/gonum/graph"
)

func Test_simpleWarehouseRobot_ID(t *testing.T) {
	type fields struct {
		id       common.RobotID
		location graph.Node
		task     common.Task
		path     []graph.Node
	}
	tests := []struct {
		name   string
		fields fields
		want   common.RobotID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &simpleWarehouseRobot{
				id:       tt.fields.id,
				location: tt.fields.location,
				task:     tt.fields.task,
				path:     tt.fields.path,
			}
			if got := r.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleWarehouseRobot.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleWarehouseRobot_Run(t *testing.T) {
	type fields struct {
		id       common.RobotID
		location graph.Node
		task     common.Task
		path     []graph.Node
	}
	type args struct {
		w  common.World
		tm common.TaskManager
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Trace
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &simpleWarehouseRobot{
				id:       tt.fields.id,
				location: tt.fields.location,
				task:     tt.fields.task,
				path:     tt.fields.path,
			}
			if got := r.Run(tt.args.w, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleWarehouseRobot.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleWarehouseRobot(t *testing.T) {
	type args struct {
		id       common.RobotID
		location graph.Node
	}
	tests := []struct {
		name string
		args args
		want common.Robot
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimpleWarehouseRobot(tt.args.id, tt.args.location); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleWarehouseRobot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeAbleToGetRobots(t *testing.T) {
	w := CreateWarehouseWorld()

	for _, r := range w.GetRobots() {
		if r == nil {
			t.Errorf("We got a problem, robot shouldn't be nil")
			t.FailNow()
		}
	}
}
func TestHasTasks(t *testing.T) {
	stm := CreateSimulatedTaskManager()
	t1 := common.NewTimePriorityTask()
	t2 := common.NewTimePriorityTask()
	t1.Origin = graph.Empty.Node()
	t2.Origin = graph.Empty.Node()
	t1.Destination = graph.Empty.Node()
	t2.Destination = graph.Empty.Node()
	stm.AddTask(t1)
	stm.AddTask(t2)
	// Expect HasTask to Be True
	if stm.HasTasks() {
		// no problem
	} else {
		t.Errorf("The HasTasks should return true, given no claim has been done")
	}
	err := stm.TaskUpdate(t1.GetTaskID(), common.Completed)
	if err != nil {
		t.Error("Upfate failed")
		t.FailNow()
	}
	// Expect HasTask to Be True
	if stm.HasTasks() {
		// no problem
	} else {
		t.Errorf("The HasTasks should return true, though 1 out of 2 is done")
	}
	err = stm.TaskUpdate(t2.GetTaskID(), common.Completed)
	if err != nil {
		t.Error("Upfate failed")
		t.FailNow()
	}
	// Expect HasTask to Be True
	if stm.HasTasks() {
		t.Errorf("The HasTasks should return false, 2 out of 2 are done")
	}
}

func TestRobotClaimTask(t *testing.T) {

	stm := CreateSimulatedTaskManager()

	w := CreateWarehouseWorld()
	t1 := common.NewTimePriorityTask()
	t1.Origin = w.GetGraph().Node(1)
	t1.Destination = w.GetGraph().Node(2)
	stm.AddTask(t1)
	id, _ := uuid.NewUUID()
	robot := NewSimpleWarehouseRobot(id, w.graph.Node(1))
	r := robot.(*simpleWarehouseRobot)
	if r.task != nil {
		t.Error("Shouldn't have any task on newly instatiated item")
	}
	r.Run(w, stm)

	if stm.HasTasks() {
		t.Errorf("Robot failed to update task status")
		t.FailNow()
	}

	if r.task == nil {
		t.Error("Failed to claim task")
	}
}
func TestRobotClaimTaskMoveAndDelivery(t *testing.T) {
	w := CreateWarehouseWorld()
	robots := w.GetRobots()
	if len(robots) <= 0 {
		t.Error("Not enough robots")
	}
	for _, i := range robots {
		if i.Location() == nil {
			t.Errorf("Robot lication is nil")
			t.FailNow()
		}
	}
	stm := CreateSimulatedTaskManager()
	t1 := common.NewTimePriorityTask()
	t2 := common.NewTimePriorityTask()
	t1.Origin = w.graph.Node(1)
	t2.Origin = w.graph.Node(1)
	t1.Destination = w.graph.Node(6)
	t2.Destination = w.graph.Node(5)
	stm.AddTask(t1)
	stm.AddTask(t2)
	// cycle to claim tasks
	trace := []common.Trace{}
	for _, i := range robots {
		trace = append(trace, i.Run(w, stm))
	}

	if len(stm.GetAllTasks()) > 1 {
		t.Errorf("Robot Failed to claim tasks")
	}
	tclaimed := false
	for _, ts := range trace {
		if ts.Target != nil && ts.Target.ID() == 6 {
			tclaimed = true
		}
	}
	if !tclaimed {
		t.Error("Failed to emit trace of t1")
	}
	tclaimed = false
	for _, ts := range trace {
		if ts.Target != nil && ts.Target.ID() == 5 {
			tclaimed = true
		}
	}
	if !tclaimed {
		t.Error("Failed to emit trace of t2")
	}
	if stm.ActiveCount() != 2 {
		t.Errorf("Failed. Acive task count should be 2, yet received %d", stm.ActiveCount())
		t.FailNow()
	}
	// cycle to move to targets
	for _, i := range robots {
		i.Run(w, stm)
	}
	for _, i := range robots {
		i.Run(w, stm)
	}
	if len(stm.GetAllTasks()) != 0 {
		t.Error("Added two basic tasks, each should take 1 cycle to finish. Yet, it still is not done")
	}
	if stm.FinishedCount() != 2 {
		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
		t.FailNow()
	}
}

func TestRobotGenerateActionPlan(t *testing.T) {
	w := CreateWarehouseWorld()
	r := w.GetRobots()[0]
	// stm := CreateSimulatedTaskManager()
	t1 := common.NewTimePriorityTask()
	t1.Origin = w.graph.Node(2)
	t1.Destination = w.graph.Node(6)
	// stm.AddTask(t1)
	if r.Location() != w.GetGraph().Node(1) {
		t.Errorf("The location of the robot initialized is incorrect")
		t.Fail()
	}
	act := PlanTaskAction(r.Location(), t1)
	if act.GetType() != action.Move {
		t.Errorf("First Generated Action sequence should be move, got %+v", act)
		t.Fail()
	}
	act = act.GetChild()
	if act.GetType() != action.StartTask {
		t.Error("Action is expected to have child Begin After Move")
		t.Fail()
	}
	act = act.GetChild()
	if act.GetType() != action.Move {
		t.Error("Action is expected to have child Move after Beging Action")
	}
	act = act.GetChild()
	if act.GetType() != action.EndTask {
		t.Error("Action is expected to have child End after move again")
	}
}

// func TestRobotCanExecuteWithMoveInSimultation(t *testing.T) {
// 	w := CreateWarehouseWorld()
// 	robots := w.GetRobots()
// 	if len(robots) <= 0 {
// 		t.Error("Not enough robots")
// 	}
// 	for _, i := range robots {
// 		if i.Location() == nil {
// 			t.Errorf("Robot lication is nil")
// 			t.FailNow()
// 		}
// 	}
// 	stm := CreateSimulatedTaskManager()
// 	t1 := common.NewTimePriorityTask()
// 	t2 := common.NewTimePriorityTask()
// 	t1.Origin = w.graph.Node(2)
// 	t2.Origin = w.graph.Node(2)
// 	t1.Destination = w.graph.Node(6)
// 	t2.Destination = w.graph.Node(5)
// 	stm.AddTask(t1)
// 	stm.AddTask(t2)
// 	// cycle to claim tasks
// 	trace := []common.Trace{}
// 	for _, i := range robots {
// 		trace = append(trace, i.Run(w, stm))
// 	}

// 	if len(stm.GetAllTasks()) > 1 {
// 		t.Errorf("Robot Failed to claim tasks")
// 	}
// 	tclaimed := false
// 	for _, ts := range trace {
// 		if ts.Target != nil && ts.Target.ID() == 6 {
// 			tclaimed = true
// 		}
// 	}
// 	if !tclaimed {
// 		t.Error("Failed to emit trace of t1")
// 	}
// 	tclaimed = false
// 	for _, ts := range trace {
// 		if ts.Target != nil && ts.Target.ID() == 5 {
// 			tclaimed = true
// 		}
// 	}
// 	if !tclaimed {
// 		t.Error("Failed to emit trace of t2")
// 	}
// 	if stm.ActiveCount() != 2 {
// 		t.Errorf("Failed. Acive task count should be 2, yet received %d", stm.ActiveCount())
// 		t.FailNow()
// 	}
// 	// cycle to move to targets
// 	for _, i := range robots {
// 		i.Run(w, stm)
// 	}
// 	for _, i := range robots {
// 		i.Run(w, stm)
// 	}
// 	if len(stm.GetAllTasks()) != 0 {
// 		t.Error("Added two basic tasks, each should take 1 cycle to finish. Yet, it still is not done")
// 	}
// 	if stm.FinishedCount() != 2 {
// 		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
// 		t.FailNow()
// 	}
// }
