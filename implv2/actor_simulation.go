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

package implv2

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

type RobotActor struct {
	robot   common.Robot
	stopSig bool
	stopped bool
}

func (state *RobotActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case InitMessage:
		state.robot.Init()

	case StartMessage:
		state.Run(ctx)
	}
}
func (state *RobotActor) Init(ctx actor.Context) {

}
func (state *RobotActor) markStopped() {
	state.stopped = true
}
func (state *RobotActor) Stop(ctx actor.Context) {
	state.stopSig = true
}
func (state *RobotActor) Run(ctx actor.Context) {
	defer state.markStopped()
	for {
		if state.stopSig {
			break
		} else {
			log.Printf("%+v", state.robot.Run())
		}
	}
}

type InitMessage struct {
}
type StartMessage struct {
}
type EndMessage struct {
}
type ShutdownMessage struct {
}

func NewSystemActor(done chan bool) *SystemActor {
	return &SystemActor{done, nil, nil, Uninitialized}
}

type SystemState int

const (
	Uninitialized SystemState = iota
	Initialized
	Running
	Stopped
)

type SystemActor struct {
	done   chan bool
	robots []*actor.PID
	w      common.World
	state  SystemState
}

func (sys *SystemActor) Receive(context actor.Context) {

	switch context.Message().(type) {
	case InitMessage:
		sys.Init(context)
	case StartMessage:
		sys.Run(context)
	case EndMessage:
		sys.Stop(context)
	case ShutdownMessage:
		sys.Shutdown(context)
	}

}
func (sys *SystemActor) Init(ctx actor.Context) {
	stm := task.CreateSimulatedTaskManagerSync()
	sys.w = world.CreateWarehouseWorldWithTaskManager(stm)
	props := actor.PropsFromProducer(sys.SpawnRobotActor)

	sys.w.AddTasks(methods.TaskGenerator(500, sys.w))
	log.Printf("Ingested %d tasks", len(sys.w.GetAllTasks()))
	for i := 0; i < 50; i++ {
		sys.robots = append(sys.robots, ctx.Spawn(props))
	}

	for _, r := range sys.robots {
		ctx.Send(r, InitMessage{})
	}
	sys.state = Initialized
}
func (sys *SystemActor) SpawnRobotActor() actor.Actor {
	return &RobotActor{robot.NewSimpleWarehouseRobot(uuid.New(), sys.w.GetGraph().Node(1), sys.w), false, false}
}
func (sys *SystemActor) Shutdown(ctx actor.Context) {
	for _, a := range sys.robots {
		ctx.Send(a, EndMessage{})
	}
	sys.done <- true
	sys.state = Stopped
}
func (sys *SystemActor) Stop(ctx actor.Context) {
	for {
		if sys.w.HasTasks() {

		} else {
			break
		}
	}
	sys.Shutdown(ctx)

}

func (sys *SystemActor) Run(ctx actor.Context) {
	switch sys.state {
	case Initialized:
		for _, r := range sys.robots {
			ctx.Send(r, StartMessage{})
		}
		sys.state = Running
	case Running:
		log.Print("System is already running")
	case Stopped:
		panic("Running an already stopped instance of the system")
	}

}
