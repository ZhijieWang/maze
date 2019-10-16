// Copyright © 2019 Zhijie (Bill) Wang <wangzhijiebill@gmail.com>
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
package impl

import (
	"maze/common"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
)

func Test_simpleWarehouseRobot_ID(t *testing.T) {
	type fields struct {
		id       common.RobotID
		location graph.Node
		task     common.Task
		path     []graph.Node
	}
	tests := []struct {
		name   string
		fields fields
		want   common.RobotID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &simpleWarehouseRobot{
				id:       tt.fields.id,
				location: tt.fields.location,
				task:     tt.fields.task,
				path:     tt.fields.path,
			}
			if got := r.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleWarehouseRobot.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleWarehouseRobot_Run(t *testing.T) {
	type fields struct {
		id       common.RobotID
		location graph.Node
		task     common.Task
		path     []graph.Node
	}
	type args struct {
		w  common.World
		tm common.TaskManager
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Trace
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &simpleWarehouseRobot{
				id:       tt.fields.id,
				location: tt.fields.location,
				task:     tt.fields.task,
				path:     tt.fields.path,
			}
			if got := r.Run(tt.args.w, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleWarehouseRobot.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSimpleWarehouseRobot(t *testing.T) {
	type args struct {
		id       common.RobotID
		location graph.Node
	}
	tests := []struct {
		name string
		args args
		want common.Robot
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimpleWarehouseRobot(tt.args.id, tt.args.location); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleWarehouseRobot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimultation(t *testing.T) {
	w := CreateWarehouseWorld()
	robots := w.GetRobots()
	if len(robots) <= 0 {
		t.Error("Not enough robots")
	}
	stm := CreateSimulatedTaskManager()
	for _, i := range robots {
		i.Run(w, stm)
	}
}
