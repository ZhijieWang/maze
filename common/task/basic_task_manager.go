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

package task

import "maze/common"

//BasicTaskManager implements a PassiveTaskManager interface, with procedure generation of tasks,
// to ensure the task queue size greater than the amount of robots
type BasicTaskManager struct {
	taskList []common.Task

	taskMap map[common.TaskID]common.Task
}

// GetTasks implements the GetTasks method from TaskManager Interface
func (tm *BasicTaskManager) GetTasks(i int) []common.Task {

	return tm.taskList
}

//ClaimTasks method implements necessary functions defined in PassiveTask managers. The method returns nil when operation was sucessful, else err.

//func (tm BasicTaskManager) ClaimTasks(TaskID, RobotID) error {
//	tm.taskListRWLock.Lock()
//	defer tm.taskListRWLock.Unlock()
//	return nil
//}

// TaskUpdate updates the status of the task, referred by taskID
func (tm *BasicTaskManager) TaskUpdate(taskID common.TaskID, status common.TaskStatus) error {

	t, err := tm.GetByID(taskID)
	if err != nil {
		// couldn't find
		return err
	}

	err = t.UpdateStatus(status)
	return err
}

func NewBasicTaskManager() *BasicTaskManager {

	tm := BasicTaskManager{}
	tm.taskMap = make(map[common.TaskID]common.Task)
	return &tm
}

// GetByID finds the task in Queue by ID
func (tM *BasicTaskManager) GetByID(taskID common.TaskID) (common.Task, error) {
	return tM.taskMap[taskID], nil

}

//Len returns the current length of the queue
func (tM *BasicTaskManager) Len() int { return len(tM.taskList) }

// Less is defined by comparing Task's Priority Function to give us the lowest based on priority
func (tM *BasicTaskManager) Less(i, j int) bool {
	return i < j
}

// Pop is predefined interface funciton in the heap interface.
// The function removes the minimum element (according to Less) from the heap and returns it. The complexity is O(log(n)) where n = h.Len(). It is equivalent to Remove(h, 0).
func (tM *BasicTaskManager) Pop() common.Task {
	n := len(tM.taskList)
	item := tM.taskList[n-1]
	delete(tM.taskMap, item.GetTaskID())
	tM.taskList = tM.taskList[0 : n-1]
	return item
}

// Push inserts the task item to the queue
func (tM *BasicTaskManager) Push(x common.Task) {
	tM.taskList = append(tM.taskList, x)
	tM.taskMap[x.GetTaskID()] = x
}

//Swap will swap elements and reblance the Task Queue

func (tM *BasicTaskManager) Swap(i, j int) {
	tM.taskList[i], tM.taskList[j] = tM.taskList[j], tM.taskList[i]
}

// AddTask insert task into the tasks manager
func (tM *BasicTaskManager) AddTask(t common.Task) bool {
	tM.Push(t)
	return true
}

// AddTasks insert tasks into the task manager
func (tm *BasicTaskManager) AddTasks(t []common.Task) bool {
	return true
}

// GetAllTasks

func (tm *BasicTaskManager) GetAllTasks() []common.Task {
	return tm.taskList
}

func (tm *BasicTaskManager) GetBroadcastInfo() interface{} {
	// currently return an instant of emptys struct (struct{})
	return struct{}{}
}

func (tm *BasicTaskManager) GetNext() common.Task {
	t := tm.Pop()
	delete(tm.taskMap, t.GetTaskID())
	return t
}

func (tm *BasicTaskManager) HasTasks() bool {
	return len(tm.taskList) > 0

}