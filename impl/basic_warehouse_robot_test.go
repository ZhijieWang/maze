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
	"maze/common"

	"reflect"
	"testing"

	"github.com/google/uuid"

	"gonum.org/v1/gonum/graph"
)

var (
	world  *WarehouseWorld
	stm    *SimulatedTaskManager
	t1     *common.TimePriorityTask
	t2     *common.TimePriorityTask
	t3     *common.TimePriorityTask
	t4     *common.TimePriorityTask
	robots []*simpleWarehouseRobot
)

func setup() {
	world = CreateWarehouseWorld()
	stm = CreateSimulatedTaskManager()

	robots = []*simpleWarehouseRobot{}
	robots = append(robots, NewSimpleWarehouseRobot(uuid.New(), world.graph.Node(1), world, stm))
	robots = append(robots, NewSimpleWarehouseRobot(uuid.New(), world.graph.Node(1), world, stm))
}
func addT1() {
	t1 = common.NewTimePriorityTask()
	t1.Origin = world.GetGraph().Node(1)
	t1.Destination = world.GetGraph().Node(2)
	stm.AddTask(t1)
}
func addT2() {
	t2 = common.NewTimePriorityTask()
	t2.Origin = world.GetGraph().Node(1)
	t2.Destination = world.GetGraph().Node(6)
	stm.AddTask(t2)
}
func addT3() {
	t3 = common.NewTimePriorityTask()
	t3.Origin = world.graph.Node(2)
	t3.Destination = world.graph.Node(6)
	stm.AddTask(t3)

}
func addT4() {
	t4 = common.NewTimePriorityTask()
	t4.Origin = world.graph.Node(2)
	t4.Destination = world.graph.Node(9)
	stm.AddTask(t4)
}

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
			if got := r.Run(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleWarehouseRobot.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleWarehouseRobot(t *testing.T) {
	type args struct {
		id       common.RobotID
		location graph.Node
		w        common.World
		tm       common.TaskManager
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
			if got := NewSimpleWarehouseRobot(tt.args.id, tt.args.location, tt.args.w, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
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
	setup()
	addT1()
	addT2()

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
	setup()
	stm.AddTask(t1)
	id, _ := uuid.NewUUID()
	r := NewSimpleWarehouseRobot(id, world.graph.Node(1), world, stm)
	if r.task != nil {
		t.Error("Shouldn't have any task on newly instatiated item")
	}
	r.Run()

	if stm.HasTasks() {
		t.Errorf("Robot failed to update task status")
		t.FailNow()
	}

	if r.task == nil {
		t.Error("Failed to claim task")
	}
}
func TestRobotClaimTaskMoveAndDelivery(t *testing.T) {

	setup()
	addT1()
	addT2()
	for _, i := range robots {
		if i.Location() == nil {
			t.Errorf("Robot location is nil")
			t.FailNow()
		}
	}
	// cycle to claim tasks

	for _, i := range robots {
		i.Run()
	}
	for _, i := range robots {
		i.Run()
	}
	if len(stm.GetAllTasks()) > 1 {
		t.Errorf("Robot Failed to claim tasks")
	}
	for _, i := range robots {
		i.Run()
	}
	for _, i := range robots {
		i.Run()
	}
	if stm.FinishedCount() != 2 {
		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
		t.FailNow()
	}
}

func TestRobotGenerateActionPlan(t *testing.T) {
	setup()
	addT3()
	r := robots[0]

	if r.Location() != world.GetGraph().Node(1) {
		t.Errorf("The location of the robot initialized is incorrect")
		t.Fail()
	}
	act := PlanTaskAction(world.GetGraph(), r.Location(), t3)
	if act.GetType() == common.ActionTypeMove && act.HasChild() && (act.(*common.MoveAction).Start == world.GetGraph().Node(1)) && (act.(*common.MoveAction).End == world.GetGraph().Node(2)) && len(act.(*common.MoveAction).Path) == 1 {

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
	if act.GetType() == common.ActionTypeMove && act.(*common.MoveAction).Start == t3.GetOrigination() && act.(*common.MoveAction).End == t3.GetDestination() {

	} else {
		t.Errorf("Action is expected to have Move after Beging Action actual is %+v", act)
	}
	act = act.GetChild()
	if act.GetType() != common.ActionTypeEndTask {
		t.Error("Action is expected to have child End after move")
	}
	act = act.GetChild()
	if act.GetType() == common.ActionTypeNull {

	} else {
		t.Errorf("Action is set to null for robot to be idle.")
	}
}
func TestRobotCanExecuteTaskPlan(t *testing.T) {
	setup()

	addT3()
	r := robots[0]

	act := PlanTaskAction(world.GetGraph(), r.Location(), t3)
	r.act = act
	r.task = t3
	node, act := r.Execute(world.GetGraph(), stm)

	if node == t3.Origin && act.GetType() == common.ActionTypeStartTask {

	} else {
		t.Errorf("target should be the task start location, actual target is %+v", node)
	}
	r.act = act
	r.location = node
	node, act = r.Execute(world.GetGraph(), stm)
	if node == t3.GetOrigination() && act.GetType() == common.ActionTypeMove && act.(*common.MoveAction).End == t3.GetDestination() {

	} else {
		t.Errorf("Failed to prepare for next step of move after begin task")
	}
}

func TestRobotCanExecuteTaskPlanMultiStep(t *testing.T) {
	setup()
	addT3()
	r := robots[0]
	trace := r.Run()

	if trace.Source == world.GetGraph().Node(1) && trace.Target == world.GetGraph().Node(2) {
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
	trace = r.Run()

	if trace.Source == t3.GetOrigination() && trace.Target == t3.GetOrigination() {
		// execute beging task action, stay at the source node
	} else {
		t.Errorf("Exepct this step to perform begin task step. Which remains at the task start position")
		t.Fail()
	}
	if r.act.GetType() == common.ActionTypeMove && r.act.(*common.MoveAction).End == t3.GetDestination() {
	} else {
		t.Errorf("Expect next step move to target location after execute beging action.\n What is the actual action? %+v", r.act)
	}
	trace = r.Run()
	if trace.Source == t3.GetOrigination() && trace.Target == t3.GetDestination() {

	} else {
		t.Errorf("Should move to final destination with execution. What actual move was %+v", trace)
	}
	if r.act.GetType() == common.ActionTypeEndTask {
		// next should end execution
	} else {
		t.Errorf("Should plan to end the task execution")
	}
	trace = r.Run()
	if trace.Source == t3.GetDestination() && trace.Target == t3.GetDestination() {
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

func TestRobotCanExecuteMoveInSimulation(t *testing.T) {

	setup()

	addT3()
	addT4()
	if len(robots) <= 0 {
		t.Error("Not enough robots")
	}
	for _, i := range robots {
		if i.Location() == nil {
			t.Errorf("Robot on is nil")
			t.FailNow()
		}
	}

	// cycle to claim tasks
	trace := []common.Trace{}
	for _, i := range robots {
		trace = append(trace, i.Run())
	}

	if len(stm.GetAllTasks()) > 1 {
		t.Errorf("Robot Failed to claim tasks")
	}
	tclaimed := false
	for _, ts := range trace {
		if ts.Target != nil && ts.Target == t3.GetOrigination() {
			tclaimed = true
		}
	}
	if !tclaimed {
		t.Error("Failed to emit trace of t1")
	}
	tclaimed = false
	for _, ts := range trace {
		if ts.Target != nil && ts.Target == t4.GetOrigination() {
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
		i.Run()
	}
	for _, i := range robots {
		i.Run()
	}
	for _, i := range robots {
		i.Run()
	}
	for _, i := range robots {
		i.Run()
	}
	if len(stm.GetAllTasks()) != 0 {
		t.Error("Added two basic tasks, each should take 1 cycle to finish. Yet, it still is not done")
	}
	if stm.FinishedCount() != 2 {
		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
		t.FailNow()
	}
}
