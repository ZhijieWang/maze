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

package test

import (
	"github.com/google/uuid"
	"maze/common"
	"maze/common/simulation"
	"maze/common/task"
	"maze/common/trace"
	"maze/common/world"
	"testing"
	"time"
)

type BasicObserver struct {
	count int
}

func (b *BasicObserver) Notify(data interface{}) {

	if data != struct{}{} {
		b.count += 1
	}
}
func (b *BasicObserver) GetChannel() chan interface{} {
	return nil
}
func TestSimulationRunResult(t *testing.T) {
	s := simulation.CreateCentralizedSimulation()
	obs := BasicObserver{}
	s.Init()

	err := s.Run(&obs)
	if err != nil {
		t.Errorf("Execution failed")
	}
	s.Stop()
	if obs.count == 0 {
		t.Errorf("Expect some run. 0 run")
	}

}

type traceObserver struct {
	traces []common.Trace
}

func (b *traceObserver) Notify(data interface{}) {
	t, ok := data.(common.Trace)
	if ok {
		b.traces = append(b.traces, t)
	}
}
func (n *traceObserver) GetChannel() chan interface{} {
	return nil
}
func TestSimulationExecuteTask(t *testing.T) {
	s := simulation.CreateCentralizedSimulation()
	obs := traceObserver{}
	s.Init()
	err := s.Run(&obs)
	if err != nil {
		t.Errorf("Execution failed")
	}
	if len(obs.traces) == 0 {
		t.Error("Failed to capture run traces")
	}
	found := false
	for _, i := range obs.traces {
		switch rTrace := i.(type) {
		case *trace.MoveTrace:

			if rTrace.RobotID != uuid.Nil {
				found = true
			}
		}
	}
	if !found {
		t.Error("No proper traces were generated from the run")
	}
}

type taskAckDistributedObserver struct {
	Receiver chan interface{}
	counter  int
	done     chan interface{}
}

func (t *taskAckDistributedObserver) GetChannel() chan interface{} {
	return t.Receiver
}

func (t *taskAckDistributedObserver) Notify(data interface{}) {
	for {

		if t.counter == 0 {
			t.done <- nil
		}
		select {

		case data = <-t.Receiver:
			if _, ok := data.(trace.TaskExecutionTrace); ok {
				t.counter--
			} else {

			}

		}
	}
}

func TestConcurrentSimulation(t *testing.T) {
	w := world.CreateWarehouseWorldWithTaskManager(task.CreateSimulatedTaskManagerSync())
	var s = simulation.System{}
	s.W = w
	s.NumBot = 5
	done := make(chan interface{})
	obs := taskAckDistributedObserver{make(chan interface{}), 10, done}
	s.Init()
	s.Start(&obs)
	go obs.Notify(nil)

	select {
	case <-done:
		s.Stop()
	// All done!
	case <-time.After(5000 * time.Millisecond):
		s.Stop()
		t.Error("Failed After 500 ms, timeout")
		t.Fail()
	}
	if len(w.GetAllTasks()) == 0 {

	} else {
		t.Fail()
	}

}
