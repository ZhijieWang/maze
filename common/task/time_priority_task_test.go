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
	"reflect"
	"testing"
	"time"

	"gonum.org/v1/gonum/graph"
)

func TestNewTimePriorityTask(t *testing.T) {
	tests := []struct {
		name string
		want *TimePriorityTask
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTimePriorityTask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimePriorityTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePriorityTask_Priority(t *testing.T) {
	type fields struct {
		ID              common.TaskID
		Origin          graph.Node
		Destination     graph.Node
		Status          common.TaskStatus
		OriginationTime time.Time
		CompletionTime  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpt := TimePriorityTask{
				ID:              tt.fields.ID,
				Origin:          tt.fields.Origin,
				Destination:     tt.fields.Destination,
				Status:          tt.fields.Status,
				OriginationTime: tt.fields.OriginationTime,
				CompletionTime:  tt.fields.CompletionTime,
			}
			if got := tpt.Priority(); got != tt.want {
				t.Errorf("TimePriorityTask.Priority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePriorityTask_GetTaskID(t *testing.T) {
	type fields struct {
		ID              common.TaskID
		Origin          graph.Node
		Destination     graph.Node
		Status          common.TaskStatus
		OriginationTime time.Time
		CompletionTime  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   common.TaskID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpt := TimePriorityTask{
				ID:              tt.fields.ID,
				Origin:          tt.fields.Origin,
				Destination:     tt.fields.Destination,
				Status:          tt.fields.Status,
				OriginationTime: tt.fields.OriginationTime,
				CompletionTime:  tt.fields.CompletionTime,
			}
			if got := tpt.GetTaskID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimePriorityTask.GetTaskID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePriorityTask_GetDestination(t *testing.T) {
	type fields struct {
		ID              common.TaskID
		Origin          graph.Node
		Destination     graph.Node
		Status          common.TaskStatus
		OriginationTime time.Time
		CompletionTime  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   graph.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpt := TimePriorityTask{
				ID:              tt.fields.ID,
				Origin:          tt.fields.Origin,
				Destination:     tt.fields.Destination,
				Status:          tt.fields.Status,
				OriginationTime: tt.fields.OriginationTime,
				CompletionTime:  tt.fields.CompletionTime,
			}
			if got := tpt.GetDestination(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimePriorityTask.GetDestination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePriorityTask_GetOrigination(t *testing.T) {
	type fields struct {
		ID              common.TaskID
		Origin          graph.Node
		Destination     graph.Node
		Status          common.TaskStatus
		OriginationTime time.Time
		CompletionTime  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   graph.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpt := TimePriorityTask{
				ID:              tt.fields.ID,
				Origin:          tt.fields.Origin,
				Destination:     tt.fields.Destination,
				Status:          tt.fields.Status,
				OriginationTime: tt.fields.OriginationTime,
				CompletionTime:  tt.fields.CompletionTime,
			}
			if got := tpt.GetOrigination(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimePriorityTask.GetOrigination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePriorityTask_UpdateStatus(t *testing.T) {
	type fields struct {
		ID              common.TaskID
		Origin          graph.Node
		Destination     graph.Node
		Status          common.TaskStatus
		OriginationTime time.Time
		CompletionTime  time.Time
	}
	type args struct {
		status common.TaskStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpt := TimePriorityTask{
				ID:              tt.fields.ID,
				Origin:          tt.fields.Origin,
				Destination:     tt.fields.Destination,
				Status:          tt.fields.Status,
				OriginationTime: tt.fields.OriginationTime,
				CompletionTime:  tt.fields.CompletionTime,
			}
			if err := tpt.UpdateStatus(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("TimePriorityTask.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
