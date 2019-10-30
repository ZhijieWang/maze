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

import (
	"maze/common"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// TimePriorityTask defines the priority of tasks via its OriginationTime
type TimePriorityTask struct {
	ID          common.TaskID
	Origin      graph.Node
	Destination graph.Node
	Status      common.TaskStatus
	//Carrier         RobotID
	OriginationTime time.Time
	CompletionTime  time.Time
}

//NewTimePriorityTask implements a basic constructor
func NewTimePriorityTask() *TimePriorityTask {
	id, _ := uuid.NewUUID()
	return &TimePriorityTask{
		ID:              id,
		Status:          common.Unassigned,
		OriginationTime: time.Now(),
	}
}

// Priority function of TimePriorityTask implements interface functions for PriorityTask
func (tpt TimePriorityTask) Priority() int64 {
	return tpt.OriginationTime.Unix()
}

// GetTaskID function of TimePriorityTask implements interface function for Task interface
func (tpt TimePriorityTask) GetTaskID() common.TaskID {
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
func (tpt TimePriorityTask) UpdateStatus(status common.TaskStatus) error {
	(&tpt).Status = status
	return nil
}

func (tpt TimePriorityTask) GetStatus() common.TaskStatus {
	return tpt.Status
}
func NewTimePriorityTaskWithParameter(start, end common.Location) *TimePriorityTask {
	id, _ := uuid.NewUUID()
	return &TimePriorityTask{
		ID:              id,
		Origin:          start,
		Destination:     end,
		Status:          common.Unassigned,
		OriginationTime: time.Now(),
	}
}
