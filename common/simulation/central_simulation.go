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

package simulation

import (
	"github.com/google/uuid"
	"log"
	"math/rand"
	"maze/common"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
)

type CentralizedSimulation struct {
	World      common.World
	TM         common.TaskManager
	Iterations int
	inited     bool
}

func CreateCentralizedSimulation() *CentralizedSimulation {
	return &CentralizedSimulation{Iterations: 10}
}

func (sim *CentralizedSimulation) Init() {

	sim.TM = task.CreateSimulatedTaskManager()
	sim.World = world.CreateWorld(sim.TM)
	l := sim.World.GetGraph().Nodes().Len()
	var numRobots = 5
	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}
		sim.World.AddRobot(robot.NewSimpleWarehouseRobot(rID, sim.World.GetGraph().Node(int64(rand.Intn(l))), sim.World))
	}

	for i := 0; i < 20; i++ {

		t := task.NewTimePriorityTask()
		t.Origin = sim.World.GetGraph().Node(int64(rand.Intn(l) + 1))
		t.Destination = sim.World.GetGraph().Node(int64(rand.Intn(l) + 1))
		if t.Origin == nil || t.Destination == nil {
			panic("Failed Initialization")
		}
		sim.World.AddTask(t)
	}
	sim.inited = true
}
func (sim *CentralizedSimulation) Run(obs common.Observer) error {
	if !sim.inited {
		panic("System enter the run mode before proper initialization")
	}
	for i := 0; i < sim.Iterations; i++ {
		for _, i := range sim.World.GetRobots() {
			trace := i.Run()
			sim.World.UpdateRobot(i)
			obs.Notify(trace)
		}
		obs.Notify(struct {
		}{})
	}
	return nil
}

func (sim *CentralizedSimulation) Stop() bool {
	return true
}
