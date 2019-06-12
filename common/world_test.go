// Copyright © 2018 Zhijie (Bill) Wang <wangzhijiebill@gmail.com>
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
package common_test

import (
	"testing"

	"maze/common"
)

func TestCanMakeWorld(t *testing.T) {
	g := common.NewSimpleWorld()
	common.
	if len(g.GetRobots()) != 0 {
		t.Errorf("Expected Empty Start, there should be no robots")
	}
	if len(g.GetTasks()) != 0 {
		t.Errorf("Expected Empty Start, there should be no tasks")
	}
	
}

func TestCanModifyRobot(t *testing.T) {
	g := common.NewSimpleWorld()
	if (g.AddRobot
	r := g.GetRobots()

	//r[0].location = r[1].location
	//if r[0].location != g.GetRobots()[0].location {
	//	t.Errorf("Expect the robots returned to be modifiable\n")
	//}
}

//func TestCanModifyTasks(t *testing.T) {
//	g := CreateWorld(3)
//	task := make([]Task, 1)
//	task = append(task, TimePriorityTask{})
//	g.SetTasks(task)
//	if len(g.GetTasks()) == 0 {
//		t.Errorf("Expect the task list to be mutable\n")
//	}
//	if g.GetTasks()[0] != task[0] {
//		t.Errorf("Expect the Task setter methog to work, but failed")
//	}
//}
