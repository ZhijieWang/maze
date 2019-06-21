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

package common

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// TaskStatus is the enumeration of task status
type TaskStatus int

// Unassigned defines task which are not assigned to any worker
// Assigned defines task that has been assigned to a worker, either in progress or not
// Completed referrs to task that are complemented and should no longer be available in task queue, but tracakle (TTL implementation subject to detail)
const (
	Unassigned = iota
	Assigned
	Completed
)

//TaskID is the alias name for a UUID, for disambiguition purpose
type TaskID = uuid.UUID

// PriorityTask interface defines what basic interface methods for tasks to be used in TaskQueue
type PriorityTask interface {
	Task
	Priority() int64
	UpdateStatus(status TaskStatus) error
}

// Task defines the data structure holding the task information
type Task interface {
	GetTaskID() TaskID
	GetDestination() graph.Node
	GetOrigination() graph.Node
}

// TimePriorityTask defines the priority of tasks via its OriginationTime
type TimePriorityTask struct {
	ID          TaskID
	Origin      graph.Node
	Destination graph.Node
	Status      TaskStatus
	//Carrier         RobotID
	OriginationTime time.Time
	CompletionTime  time.Time
}

//NewTimePriorityTask implements a basic constructor
func NewTimePriorityTask() *TimePriorityTask {
	id, _ := uuid.NewUUID()
	return &TimePriorityTask{
		ID:              id,
		Status:          Unassigned,
		OriginationTime: time.Now(),
	}

}

// Priority function of TimePriorityTask implements interface functions for PriorityTask
func (tpt TimePriorityTask) Priority() int64 {
	return tpt.OriginationTime.Unix()
}

// GetTaskID function of TimePriorityTask implements interface function for Task inteface
func (tpt TimePriorityTask) GetTaskID() TaskID {
	return tpt.ID
}

// GetDestination function of TimePriorityTask implements interface functions for Task interface
func (tpt TimePriorityTask) GetDestination() graph.Node {
	return tpt.Destination
}

// GetOrigination function of TimePriorityTask implements interface function for Task interface
func (tpt TimePriorityTask) GetOrigination() graph.Node {
	return tpt.Origin
}

// UpdateStatus update the current status of the task, further logic is needed
func (tpt TimePriorityTask) UpdateStatus(status TaskStatus) error {
	tpt.Status = status
	return nil
}

//TaskManager defines task manager interfaces. All task generator, coordinator must follow this type
type TaskManager interface {
	Run()
	GetBoradcastInfo() interface{}
	GetTasks() []PriorityTask
	TaskUpdate(taskID TaskID, status TaskStatus) error
}

// PassiveTaskManager extends the TaskManager interface and allows robots to claim tasks
type PassiveTaskManager interface {
	TaskManager
	//	ClaimTask(Task, RobotID) error
}

//BasicTaskManager implements a PassiveTaskManager interface, with procedure generation of tasks,
// to ensure the task queue size greater than the amount of robots
type BasicTaskManager struct {
	taskBufferSize int
	taskList       TaskQueue
	taskListRWLock *sync.RWMutex
	taskArchive    []Task
}

// GetTasks implements the GetTasks method from TaskManager Interface
func (tm BasicTaskManager) GetTasks() []PriorityTask {
	tm.taskListRWLock.RLock()
	defer tm.taskListRWLock.RUnlock()
	return tm.taskList.taskList
}

//ClaimTasks method implements necessary functions defined in PassiveTask managers. The method returns nil when operation was sucessful, else err.

//func (tm BasicTaskManager) ClaimTasks(TaskID, RobotID) error {
//	tm.taskListRWLock.Lock()
//	defer tm.taskListRWLock.Unlock()
//	return nil
//}

// TaskUpdate updates the status of the task, referred by taskID
func (tm BasicTaskManager) TaskUpdate(taskID TaskID, status TaskStatus) error {
	err := (*tm.taskList.GetByID(taskID)).UpdateStatus(status)
	if err != nil {
		return nil
	}
	return nil
}

//TaskQueue implements a FIFO heap based task queue while allow task to be looked up by TaskID. TaskQueue currently is an approximation for 2 way indexed database.
type TaskQueue struct {
	taskList []PriorityTask
	taskMap  map[TaskID]*PriorityTask
}

// GetByID finds the task in Queue by ID
func (tq TaskQueue) GetByID(taskID TaskID) *PriorityTask {
	return tq.taskMap[taskID]

}

//Len returns the current lenght of the taskqueue object
func (tq TaskQueue) Len() int { return len(tq.taskList) }

// Less is defined by comparing Task's Priority Function to give us the lowest based on priority
func (tq TaskQueue) Less(i, j int) bool {
	return tq.taskList[i].Priority() < tq.taskList[j].Priority()
}

// Pop is predefined interface funciton in the heap interface.
// The function removes the minimum element (according to Less) from the heap and returns it. The complexity is O(log(n)) where n = h.Len(). It is equivalent to Remove(h, 0).
func (tq *TaskQueue) Pop() PriorityTask {
	n := len(tq.taskList)
	item := tq.taskList[n-1]
	delete(tq.taskMap, item.GetTaskID())
	tq.taskList = tq.taskList[0 : n-1]
	return item
}

// Push inserts the task item to the queue
func (tq *TaskQueue) Push(x PriorityTask) {
	tq.taskList = append(tq.taskList, x)
	tq.taskMap[x.GetTaskID()] = &x
}

//Swap will swap elements and reblance the Task Queue

func (tq TaskQueue) Swap(i, j int) {
	tq.taskList[i], tq.taskList[j] = tq.taskList[j], tq.taskList[i]
}

// NewTaskQueue is the constructor method for TaskQueue to initialize necessary values.
func NewTaskQueue() *TaskQueue {
	tq := TaskQueue{}
	tq.taskMap = make(map[TaskID]*PriorityTask)
	return &tq
}
