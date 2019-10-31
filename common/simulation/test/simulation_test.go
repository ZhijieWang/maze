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

package test_test

import (
	"maze/common/simulation"
	"testing"
)

type BasicObserver struct {
	count int
}

func (b *BasicObserver) OnNotify(data interface{}) {

	if data != struct{}{} {
		b.count += 1
	}
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

//type traceObserver struct {
//	traces []common.Trace
//}
//
//func (b *traceObserver) OnNotify(data interface{}) {
//	t, ok := data.(common.Trace)
//	if ok {
//		b.traces = append(b.traces, t)
//	}
//}
//
//
//func TestSimulationExecuteTask(t *testing.T) {
//	s := simulation.CreateCentralizedSimulation()
//	obs := traceObserver{}
//	err := s.Run(&obs)
//	if err != nil {
//		t.Errorf("Execution failed")
//	}
//	if len(obs.traces) == 0 {
//		t.Error("Failed to capture run traces")
//	}
//	found := false
//	for _, i := range obs.traces {
//		if i.(*trace.MoveTrace).RobotID != uuid.Nil {
//			found = true
//		}
//	}
//	if !found {
//		t.Error("No proper traces were generated from the run")
//	}
//}
//
//type taskAckDistributedObserver struct{
//	counter int
//}
//func (t *taskAckDistributedObserver)OnNotify(data interface{}){
//	if _, ok:= data.(trace.TaskExecutionTrace); ok{
//		t.counter --
//	}
//}
//
//func TestConcurrentSimulation( t *testing.T){
//	var s = simulation.System{}
//	s.NumBot = 5
//	obs := taskAckDistributedObserver{10}
//	s.Init()
//	s.Start(&obs)
//	done := make(chan struct{})
//	go func(){
//		for {
//			if obs.counter ==0{
//				s.Stop()
//			}
//			close(done)
//		}
//	}()
//
//	select {
//	case <-done:
//	// All done!
//	case <-time.After(500 * time.Millisecond):
//		t.Error("Failed After 500 ms, timeout")
//		t.Fail()
//	}
//
//}
