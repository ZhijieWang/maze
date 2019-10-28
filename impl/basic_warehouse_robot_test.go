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

	for _, i := range robots {
		i.Run(w, stm)
	}
	for _, i := range robots {
		i.Run(w, stm)
	}
	if len(stm.GetAllTasks()) > 1 {
		t.Errorf("Robot Failed to claim tasks")
	}
	for _, i := range robots {
		i.Run(w, stm)
	}
	for _, i := range robots {
		i.Run(w, stm)
	}
	if stm.FinishedCount() != 2 {
		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
		t.FailNow()
	}
}

func TestRobotGenerateActionPlan(t *testing.T) {
	w := CreateWarehouseWorld()
	r := w.GetRobots()[0]

	t1 := common.NewTimePriorityTask()
	t1.Origin = w.graph.Node(2)
	t1.Destination = w.graph.Node(6)

	if r.Location() != w.GetGraph().Node(1) {
		t.Errorf("The location of the robot initialized is incorrect")
		t.Fail()
	}
	act := PlanTaskAction(w.GetGraph(), r.Location(), t1)
	if act.GetType() == common.ActionTypeMove && act.HasChild() && (act.(*common.MoveAction).Start == w.GetGraph().Node(1)) && (act.(*common.MoveAction).End == w.GetGraph().Node(2)) && len(act.(*common.MoveAction).Path) == 1 {

	} else {
		t.Errorf("First Generated Action sequence should be ActionTypeMove, got %+v", act)
		t.Fail()
	}
	act = act.GetChild()
	if act.GetType() != common.ActionTypeStartTask {
		t.Errorf("Action is expected to have child BeginTask, actual is %+v", act)
		t.Fail()
	}
	act = act.GetChild()
	if act.GetType() == common.ActionTypeMove && act.(*common.MoveAction).Start == t1.GetOrigination() && act.(*common.MoveAction).End == t1.GetDestination() {

	} else {
		t.Errorf("Action is expected to have Move after Beging Action actual is %+v", act)
	}
	act = act.GetChild()
	if act.GetType() != common.ActionTypeEndTask {
		t.Error("Action is expected to have child End aften.ActionTypeMove again")
	}
	act = act.GetChild()
	if act.GetType() == common.ActionTypeNull {

	} else {
		t.Errorf("Action is set to null for robot to be idle.")
	}
}
func TestRobotCanExecuteTaskPlan(t *testing.T) {
	w := CreateWarehouseWorld()
	r := w.GetRobots()[0]
	stm := CreateSimulatedTaskManager()
	t1 := common.NewTimePriorityTask()
	t1.Origin = w.graph.Node(2)
	t1.Destination = w.graph.Node(6)
	stm.AddTask(t1)
	if r.Location() != w.GetGraph().Node(1) {
		t.Errorf("The location of the robot initialized is incorrect")
		t.Fail()
	}
	// targetAct := action.CreateMoveAction(w.GetGraph().Node(1), w.GetGraph().Node(2))

	act := PlanTaskAction(w.GetGraph(), r.Location(), t1)
	r.(*simpleWarehouseRobot).act = act
	r.(*simpleWarehouseRobot).task = t1
	node, act := r.(*simpleWarehouseRobot).Execute(w.GetGraph(), stm)

	if node == t1.Origin && act.GetType() == common.ActionTypeStartTask {

	} else {
		t.Errorf("target should be the task start location, actual target is %+v", node)
	}
	r.(*simpleWarehouseRobot).act = act
	r.(*simpleWarehouseRobot).location = node
	node, act = r.(*simpleWarehouseRobot).Execute(w.GetGraph(), stm)
	if node == t1.GetOrigination() && act.GetType() == common.ActionTypeMove && act.(*common.MoveAction).End == t1.GetDestination() {

	} else {
		t.Errorf("Failed to prepare for next step of move after begin task")
	}
}

func TestRobotCanExecuteTaskPlanMultiStep(t *testing.T) {
	w := CreateWarehouseWorld()
	r := w.GetRobots()[0].(*simpleWarehouseRobot)
	stm := CreateSimulatedTaskManager()
	t1 := common.NewTimePriorityTask()
	t1.Origin = w.graph.Node(2)
	t1.Destination = w.graph.Node(6)
	stm.AddTask(t1)
	trace := r.Run(w, stm)

	if trace.Source == w.GetGraph().Node(1) && trace.Target == w.GetGraph().Node(2) {
		// Move one step
	} else {
		t.Errorf("First step should be moving from 1 to 2, actual trace is %+v", trace)
		t.Fail()
	}
	if r.act.GetType() == common.ActionTypeStartTask {
		// a Move, the action next should be beging task
	} else {
		t.Errorf("After 1 step move, the next pending task should be begin task, but actual is %v : %+v", r.act.GetType(), r.act)
		t.Fail()
	}
	trace = r.Run(w, stm)

	if trace.Source == t1.GetOrigination() && trace.Target == t1.GetOrigination() {
		// execute beging task action, stay at the source node
	} else {
		t.Errorf("Exepct this step to perform begin task step. Which remains at the task start position")
		t.Fail()
	}
	if r.act.GetType() == common.ActionTypeMove && r.act.(*common.MoveAction).End == t1.GetDestination() {
	} else {
		t.Errorf("Expect next step move to target location after execute beging action.\n What is the actual action? %+v", r.act)
	}
	trace = r.Run(w, stm)
	if trace.Source == t1.GetOrigination() && trace.Target == t1.GetDestination() {

	} else {
		t.Errorf("Should move to final destination with execution. What actual move was %+v", trace)
	}
	if r.act.GetType() == common.ActionTypeEndTask {
		// next should end execution
	} else {
		t.Errorf("Should plan to end the task execution")
	}
	trace = r.Run(w, stm)
	if trace.Source == t1.GetDestination() && trace.Target == t1.GetDestination() {
		// wrap up task
	} else {
		t.Errorf("Execte end task ")
	}
	if r.act == common.Null() {
		// Should be nothing left
	} else {
		t.Errorf("Should release the robot to idle. Actual: %+v", r.act)
	}
}

func TestRobotCanExecuteMoveInSimultation(t *testing.T) {
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
	t1.Origin = w.graph.Node(2)
	t2.Origin = w.graph.Node(2)
	t1.Destination = w.graph.Node(6)
	t2.Destination = w.graph.Node(9)
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
		if ts.Target != nil && ts.Target == t1.GetOrigination() {
			tclaimed = true
		}
	}
	if !tclaimed {
		t.Error("Failed to emit trace of t1")
	}
	tclaimed = false
	for _, ts := range trace {
		if ts.Target != nil && ts.Target == t2.GetOrigination() {
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
	// cycle tn.ActionTypeMove to targets
	for _, i := range robots {
		i.Run(w, stm)
	}
	for _, i := range robots {
		i.Run(w, stm)
	}
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
