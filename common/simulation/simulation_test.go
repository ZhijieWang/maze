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

package simulation_test

import (
	"log"
	"maze/common"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
	"testing"

	"github.com/google/uuid"
)

type centralizedSimulation struct {
	world common.World
	tm    common.TaskManager
}

func CreateCentralizedSimulation() common.Simulation {

	var c = centralizedSimulation{}
	c.tm = task.NewBasicTaskManager()
	c.world = world.CreateWorld(1, task.NewBasicTaskManager())
	// c.world = common.CreateBlankWorld()
	var numRobots int = 5
	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}
		c.world.AddRobot(robot.NewSimpleRobot(rID, c.world.GetGraph().Nodes().Node(), c.world, c.tm))
	}
	for i := 0; i < 20; i++ {
		t := task.NewTimePriorityTask()
		t.Destination = c.world.GetGraph().Nodes().Node()
		c.tm.AddTask(t)
	}

	return c
}

func (sim centralizedSimulation) Init() {

}
func (sim centralizedSimulation) Run(obs common.Observer) error {
	for _, i := range sim.world.GetRobots() {
		trace := i.Run()
		sim.world.UpdateRobot(i)
		obs.OnNotify(trace)
	}
	obs.OnNotify(nil)

	return nil
}

func (sim centralizedSimulation) Stop() bool {
	return true
}

type basicObserver struct {
	count int
}

func (b *basicObserver) OnNotify(data interface{}) {

	if data != struct{}{} {
		b.count += 1
	}
}

type traceObserver struct {
	traces []common.Trace
}

func (b *traceObserver) OnNotify(data interface{}) {
	t, ok := data.(common.Trace)
	if ok {
		b.traces = append(b.traces, t)
	}
}

func TestSimulationRunResult(t *testing.T) {
	s := CreateCentralizedSimulation()
	obs := basicObserver{}
	err := s.Run(&obs)
	if err != nil {
		t.Errorf("Execution failed")
	}
	s.Stop()
	if obs.count == 0 {
		t.Errorf("Expect some run. 0 run")
	}

}

func TestSimulationExecuteTask(t *testing.T) {
	s := CreateCentralizedSimulation()
	obs := traceObserver{}
	err := s.Run(&obs)
	if err != nil {
		t.Errorf("Execution failed")
	}
	if len(obs.traces) == 0 {
		t.Error("Failed to capture run traces")
	}
	found := false
	for _, i := range obs.traces {
		if i.RobotID != uuid.Nil {
			found = true
		}
	}
	if !found {
		t.Error("No proper traces were generated from the run")
	}
}
