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

package impl

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"log"
	"maze/common"
	"maze/common/methods"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
)

type ShutdownMessageV1 struct {
}
type RobotActorV1 struct {
	robot   common.Robot
	stopSig bool
	stopped bool
}

func (state *RobotActorV1) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case InitMessageV1:
		state.robot.Init()

	case StartMessageV1:
		state.Run(ctx)
	}
}
func (state *RobotActorV1) Init(ctx actor.Context) {

}
func (state *RobotActorV1) markStopped() {
	state.stopped = true
}
func (state *RobotActorV1) Stop(ctx actor.Context) {
	state.stopSig = true
}
func (state *RobotActorV1) Run(ctx actor.Context) {
	defer state.markStopped()
	for {
		if state.stopSig {
			break
		} else {
			log.Printf("%+v", state.robot.Run())
		}
	}
}

func NewSystemActorV2(done chan bool) *SystemActorV2 {
	return &SystemActorV2{done, nil, nil, Uninitialized}
}

type SystemState int

const (
	Uninitialized SystemState = iota
	Initialized
	Running
	Stopped
)

type SystemActorV2 struct {
	done   chan bool
	robots []*actor.PID
	w      common.World
	state  SystemState
}

func (sys *SystemActorV2) Receive(context actor.Context) {

	switch context.Message().(type) {
	case InitMessageV1:
		sys.Init(context)
	case StartMessageV1:
		sys.Run(context)
	case EndMessageV1:
		sys.Stop(context)
	case ShutdownMessageV1:
		sys.Shutdown(context)
	}

}
func (sys *SystemActorV2) Init(ctx actor.Context) {
	stm := task.CreateSimulatedTaskManagerSync()
	sys.w = world.CreateWarehouseWorldWithTaskManager(stm)
	props := actor.PropsFromProducer(sys.SpawnRobotActor)

	sys.w.AddTasks(methods.TaskGenerator(500, sys.w))
	log.Printf("Ingested %d tasks", len(sys.w.GetAllTasks()))
	for i := 0; i < 50; i++ {
		sys.robots = append(sys.robots, ctx.Spawn(props))
	}

	for _, r := range sys.robots {
		ctx.Send(r, InitMessageV1{})
	}
	sys.state = Initialized
}
func (sys *SystemActorV2) SpawnRobotActor() actor.Actor {
	return &RobotActorV1{robot.NewSimpleWarehouseRobot(uuid.New(), sys.w.GetGraph().Node(1), sys.w), false, false}
}
func (sys *SystemActorV2) Shutdown(ctx actor.Context) {
	for _, a := range sys.robots {
		ctx.Send(a, EndMessageV1{})
	}
	sys.done <- true
	sys.state = Stopped
}
func (sys *SystemActorV2) Stop(ctx actor.Context) {
	for {
		if sys.w.HasTasks() {

		} else {
			break
		}
	}
	sys.Shutdown(ctx)

}

func (sys *SystemActorV2) Run(ctx actor.Context) {
	switch sys.state {
	case Initialized:
		for _, r := range sys.robots {
			ctx.Send(r, StartMessageV1{})
		}
		sys.state = Running
	case Running:
		log.Print("System is already running")
	case Stopped:
		panic("Running an already stopped instance of the system")
	}

}
