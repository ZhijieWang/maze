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

package trace

import (
	"gonum.org/v1/gonum/graph"
	"maze/common"
)

type MoveTrace struct {
	RobotID   common.RobotID
	Source    graph.Node
	Target    graph.Node
	Timestamp int
}

func (m *MoveTrace) GetType() common.TraceType {
	return common.MoveTraceType
}

func (m *MoveTrace) GetContent() interface{} {
	return m
}

type TaskExecutionTrace struct {
	Status  int
	TaskID  common.TaskID
	RobotID common.RobotID
}

type TaskNullActionTrace struct {
}

var TaskNullActionTraceType common.TraceType = 2

func (m TaskNullActionTrace) GetType() common.TraceType {
	return TaskNullActionTraceType
}
func (m TaskNullActionTrace) GetContent() interface{} {
	return m
}

var TaskExecutionTraceType common.TraceType = 1

func (m TaskExecutionTrace) GetType() common.TraceType {
	return TaskExecutionTraceType
}
func (m TaskExecutionTrace) GetContent() interface{} {
	return m
}
