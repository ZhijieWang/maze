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

package implv1

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"log"
	"maze/common"
	"maze/common/robot"
	"maze/common/world"
)

type InitMessage struct {
}
type StartMessage struct {
}
type EndMessage struct {
}

func NewSystemActor(s chan bool) *SystemActor {
	return &SystemActor{nil, nil, nil, s}
}

type SystemActor struct {
	w             common.World
	robots        []common.Robot
	taskGenerator func(maxTasks int, w common.World) []common.Task
	s             chan bool
}

func (sys *SystemActor) Init() {
	log.Print("System Initialized")
	sys.w = world.CreateWarehouseWorld()
	for i := 0; i < 10; i++ {
		id, _ := uuid.NewUUID()
		sys.robots = append(sys.robots, robot.NewSimpleWarehouseRobot(id, sys.w.GetGraph().Node(1), sys.w))
	}
}
func (sys *SystemActor) Run() {
	for _, a := range sys.robots {
		log.Printf("%+v ", a.Run())
	}
}
func (sys *SystemActor) Stop() {
	sys.s <- true
}
func (sys *SystemActor) Receive(context actor.Context) {
	log.Print(context.Message())
	switch context.Message().(type) {
	case InitMessage:

		sys.Init()
	case StartMessage:
		sys.Run()
	case EndMessage:
		sys.Stop()
	}
}
