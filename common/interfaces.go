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
	Equal(Action) bool
}

type DurationAction interface {
	GetStatus() ActionStatus
	SetStatus(ActionStatus)
}

// RobotID is an alias to UUID for disambiguation purpose
type RobotID = uuid.UUID
type Robot interface {
	ID() RobotID
	Init() bool
	Run() Trace
	Location() graph.Node
	Plan()
	Execute() Trace
	GetStatus() (Action, Task)
}

type TraceType int

const (
	MoveTraceType TraceType = iota
	TaskTraceType
)

// Trace is data structure to hold data of robot movement and actions, for tracking
type Trace interface {
	GetType() TraceType
	GetContent() interface{}
}

// TaskStatus is the enumeration of task status
type TaskStatus int

// Unassigned defines task which are not assigned to any worker
// Assigned defines task that has been assigned to a worker, either in progress or not
// Completed refers to task that are complemented and should no longer be available in task queue, but traceable (TTL implementation subject to detail)
const (
	Unassigned = iota
	Assigned
	Completed
)

//TaskID is the alias name for a UUID, for disambiguation purpose
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
	GetNextTask() Task
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
	TaskManager
	GetGraph() graph.Graph
	GetRobots() []Robot
	UpdateRobot(Robot) bool
	AddRobot(r Robot) bool
	ClaimTask(tid TaskID, rid RobotID) (success bool, err error)
}
type Observer interface {
	Notify(data interface{})
	GetChannel() chan interface{}
}

type Event interface {
}
type Notifier interface {
	Register(Observer)
	Deregister(Observer)
	Notify(Event)
}
type Simulation interface {
	Init()
	Run(obs Observer) error
	Stop() bool
}
type Actor interface {
	Init()
	Run(observer Observer)
	Stop()
}
