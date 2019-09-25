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
	"maze/common"
	"testing"
)

func TestCanMakeWorld(t *testing.T) {
	g := common.CreateWorld(2, common.NewBasicTaskManager())
	if len(g.GetRobots()) != 0 {
		t.Errorf("Expected Empty Start, there should be no robots")
	}
	if len(g.GetTasks()) != 0 {
		t.Errorf("Expected Empty Start, there should be no tasks")
	}

}

//func TestCanAssignRobot(t *testing.T) {
//	g := common.NewSimpleWorld()
//	var robots [5]common.Robot = [5]common.Robot{common.Robot{}}
//	g.AddRobots(robots[:])
//	r := g.GetRobots()
//	if len(r) != 5 {
//		t.Errorf("We have a problem. Expected length 5, actual length %d", len(r))
//	}
//}

func TestCanModifyTasks(t *testing.T) {
	g := common.CreateWorld(2, common.NewBasicTaskManager())
	task := make([]common.Task, 0)
	task = append(task, common.TimePriorityTask{})
	g.SetTasks(task)
	if len(g.GetTasks()) == 0 {
		t.Errorf("Expect the task list to be mutable\n")
	}
	if g.GetTasks()[0] != task[0] {
		t.Errorf("Expect the Task setter methog to work, but failed")
	}
}
