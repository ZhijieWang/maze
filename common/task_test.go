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

package common_test

import (
	"maze/common"
	"testing"
)

func TestTaskManagerPushAtomicSuccess(t *testing.T) {
	var tq common.TaskManager = common.NewBasicTaskManager()
	tq.AddTask(common.NewTimePriorityTask())

	if len(tq.GetAllTasks()) != 1 {
		t.Errorf("Insert one task into queue, expect queue size to be 1\n, current length is %d", len(tq.GetAllTasks()))
	}
}

func TestTaskManagerPushMaintainOrder(t *testing.T) {
	tq := common.NewBasicTaskManager()
	t1 := common.NewTimePriorityTask()
	t2 := common.NewTimePriorityTask()
	tq.Push(t2)
	tq.Push(t1)
	if t1 == t2 {
		t.Error("Input should be different\n")
	}
	if t1 != tq.Pop() {
		t.Errorf("Expect the task queue to maintain time order for out of order push\n")
	}
	if t2 != tq.Pop() {
		t.Errorf("Expect the task queue to maintain time order for out of order push\n")
	}

}
