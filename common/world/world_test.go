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
package world_test

import (
	"maze/common"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
	"testing"

	"github.com/google/uuid"
)

func TestCanMakeWorld(t *testing.T) {
	g := world.CreateWorld(0, task.NewBasicTaskManager())
	if len(g.GetRobots()) != 0 {
		t.Errorf("Expected Empty Start, there should be no robots")
	}
	if len(g.GetTasks()) != 0 {
		t.Errorf("Expected Empty Start, there should be no tasks")
	}

}

func TestCanAssignRobot(t *testing.T) {
	g := world.CreateWorld(2, task.NewBasicTaskManager())
	numBots := len(g.GetRobots())
	id, _ := uuid.NewUUID()
	r := robot.NewSimpleRobot(id, g.GetGraph().Node(0), g, task.NewBasicTaskManager())
	g.AddRobot(r)
	rs := g.GetRobots()
	if len(rs) != (numBots + 1) {
		t.Errorf("We have a problem. Expected length %d , actual length %d", numBots+1, len(rs))
	}
}

func TestCanModifyTasks(t *testing.T) {
	g := world.CreateWorld(2, task.NewBasicTaskManager())
	tasks := make([]common.Task, 0)
	tasks = append(tasks, task.TimePriorityTask{})
	g.SetTasks(tasks)
	if len(g.GetTasks()) == 0 {
		t.Errorf("Expect the task list to be mutable\n")
	}
	if g.GetTasks()[0] != tasks[0] {
		t.Errorf("Expect the Task setter method to work, but failed")
	}
}