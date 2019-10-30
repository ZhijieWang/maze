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
package robot_test

import (
	"maze/common"
	"maze/common/action"
	"maze/common/methods"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"

	"testing"

	"github.com/google/uuid"
)

var (
	w      common.World
	stm    *task.SimulatedTaskManager
	t1     *task.TimePriorityTask
	t2     *task.TimePriorityTask
	t3     *task.TimePriorityTask
	t4     *task.TimePriorityTask
	robots []common.Robot
)

func setup() {
	w = world.CreateWarehouseWorld()
	stm = task.CreateSimulatedTaskManager()

	robots = []common.Robot{}
	robots = append(robots, robot.NewSimpleWarehouseRobot(uuid.New(), w.GetGraph().Node(1), w, stm))
	robots = append(robots, robot.NewSimpleWarehouseRobot(uuid.New(), w.GetGraph().Node(1), w, stm))
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
func addT3() {
	t3 = task.NewTimePriorityTask()
	t3.Origin = w.GetGraph().Node(2)
	t3.Destination = w.GetGraph().Node(6)
	stm.AddTask(t3)

}
func addT4() {
	t4 = task.NewTimePriorityTask()
	t4.Origin = w.GetGraph().Node(2)
	t4.Destination = w.GetGraph().Node(9)
	stm.AddTask(t4)
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

	if r.Location() != w.GetGraph().Node(1) {
		t.Errorf("The location of the robot initialized is incorrect")
		t.Fail()
	}
	act := methods.PlanTaskAction(w.GetGraph(), r.Location(), t3)
	if act.GetType() == common.ActionTypeMove && act.HasChild() && (act.(*action.MoveAction).Start == w.GetGraph().Node(1)) && (act.(*action.MoveAction).End == w.GetGraph().Node(2)) && len(act.(*action.MoveAction).Path) == 1 {

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
	if act.GetType() == common.ActionTypeMove && act.(*action.MoveAction).Start == t3.GetOrigination() && act.(*action.MoveAction).End == t3.GetDestination() {

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

	act := methods.PlanTaskAction(w.GetGraph(), r.Location(), t3)
	r.Plan()
	//r.Execute()
	rAct, rTask := r.GetStatus()
	if rAct.Equal(act) {

	} else {
		t.Errorf("Robot Plan picked up wrong act, expecting %+v, actual %+v", act, rAct)
		t.Fail()
	}
	if rTask == t3 {

	} else {
		t.Errorf("Robot Plan picked up wrong task, expecting %+v, actual %+v", t3, rTask)
		t.Fail()
	}
	r.Execute()
	act = act.GetChild()
	rAct, rTask = r.GetStatus()
	if rAct.Equal(act) {

	} else {
		t.Errorf("Robot Execute picked up wrong act, expecting %+v, actual %+v", act, rAct)
		t.Fail()
	}
	if rTask == t3 {

	} else {
		t.Errorf("Robot Execute picked up wrong task, expecting %+v, actual %+v", t3, rTask)
		t.Fail()
	}

}

func TestRobotCanExecuteTaskPlanMultiStep(t *testing.T) {
	setup()
	addT3()
	r := robots[0]
	trace := r.Run()
	rAct, _ := r.GetStatus()
	if trace.Source == w.GetGraph().Node(1) && trace.Target == w.GetGraph().Node(2) {
		// Move one step
	} else {
		t.Errorf("First step should be moving from 1 to 2, actual trace is %+v", trace)
		t.Fail()
	}

	if rAct.GetType() == common.ActionTypeStartTask {
		// a Move, the action next should be beging task
	} else {
		t.Errorf("After 1 step move, the next pending task should be begin task, but actual is %v : %+v", rAct.GetType(), rAct)
		t.Fail()
	}
	trace = r.Run()
	rAct, _ = r.GetStatus()
	if trace.Source == t3.GetOrigination() && trace.Target == t3.GetOrigination() {
		// execute beging task action, stay at the source node
	} else {
		t.Errorf("Exepct this step to perform begin task step. Which remains at the task start position")
		t.Fail()
	}
	if rAct.GetType() == common.ActionTypeMove && rAct.(*action.MoveAction).End == t3.GetDestination() {
	} else {
		t.Errorf("Expect next step move to target location after execute beging action.\n What is the actual action? %+v", rAct)
	}
	trace = r.Run()
	rAct, _ = r.GetStatus()
	if trace.Source == t3.GetOrigination() && trace.Target == t3.GetDestination() {

	} else {
		t.Errorf("Should move to final destination with execution. What actual move was %+v", trace)
	}
	if rAct.GetType() == common.ActionTypeEndTask {
		// next should end execution
	} else {
		t.Errorf("Should plan to end the task execution")
	}
	trace = r.Run()
	rAct, _ = r.GetStatus()
	if trace.Source == t3.GetDestination() && trace.Target == t3.GetDestination() {
		// wrap up task
	} else {
		t.Errorf("Execte end task ")
	}
	if rAct == action.Null() {
		// Should be nothing left
	} else {
		t.Errorf("Should release the robot to idle. Actual: %+v", rAct)
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
	var trace []common.Trace
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
	if stm.FinishedCount() != 2 || stm.ActiveCount() != 0 {
		t.Errorf("Failed. Finished task count should be 2, yet received %d", stm.FinishedCount())
		t.FailNow()
	}
}
