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

package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

type ActionType int

const (
	PendingStatus = iota
	ActiveStatus
	EndStatus
)

type ActionStatus int

const (
	ActionTypeMove = iota
	ActionTypeStartTask
	ActionTypeEndTask
	ActionTypeNull
)

type Action interface {
	GetType() ActionType
	GetContent() interface{}
	HasChild() bool
	GetChild() Action
	SetChild(Action)
}

type DurationAction interface {
	GetStatus() ActionStatus
	SetStatus(ActionStatus)
}

// RobotID is an alias to UUID for disambiguition purpose
type RobotID = uuid.UUID
type Robot interface {
	ID() RobotID
	Init() bool
	Run() Trace
	Location() graph.Node
	Plan(graph.Graph)
	Execute(graph.Graph, TaskManager) (graph.Node, Action)
}

// Trace is data structure to hold data that can be used for path planning
type Trace struct {
	RobotID   RobotID
	Source    graph.Node
	Target    graph.Node
	Timestamp int
}

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

// PriorityTask interface defines what basic interface methods for tasks to be used in TaskQueue. The Priority accessor can be used as a comparator.
type PriorityTask interface {
	Task
	Priority() int64
}

// Task defines the data structure holding the task information
type Task interface {
	GetTaskID() TaskID
	GetDestination() graph.Node
	GetOrigination() graph.Node
	UpdateStatus(status TaskStatus) error
	GetStatus() TaskStatus
}

//TaskManager defines task manager interfaces. All task generator, coordinator must follow this type
type TaskManager interface {
	GetBroadcastInfo() interface{}
	GetAllTasks() []Task
	GetNext() Task
	GetTasks(n int) []Task
	TaskUpdate(taskID TaskID, status TaskStatus) error
	AddTask(t Task) bool
	AddTasks(tList []Task) bool
	HasTasks() bool
}

// PassiveTaskManager extends the TaskManager interface and allows robots to claim tasks
type PassiveTaskManager interface {
	TaskManager
	//	ClaimTask(Task, RobotID) error
}
type Location graph.Node

// World interface defines the behavior of World simulation
type World interface {
	GetGraph() graph.Graph
	GetRobots() []Robot
	UpdateRobot(Robot) bool
	GetTasks() []Task
	AddRobot(r Robot) bool
	SetTasks(tasks []Task) bool
	ClaimTask(tid TaskID, rid RobotID)
}
