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

type centralizedSimulation struct {
	world common.World
	tm    common.TaskManager
}

func CreateCentralizedSimulation() common.Simulation {

	var c = centralizedSimulation{}
	c.tm = task.NewBasicTaskManager()
	c.world = world.CreateWorld(c.tm)
	l := c.world.GetGraph().Nodes().Len()
	var numRobots = 5
	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}
		c.world.AddRobot(robot.NewSimpleWarehouseRobot(rID, c.world.GetGraph().Node(int64(rand.Intn(l))), c.world))
	}

	for i := 0; i < 20; i++ {

		t := task.NewTimePriorityTask()
		t.Origin = c.world.GetGraph().Node(int64(rand.Intn(l)))
		t.Destination = c.world.GetGraph().Node(int64(rand.Intn(l)))
		c.world.AddTask(t)
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
