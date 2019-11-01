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
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
)

type RobotActorV2 struct {
	robot   common.Robot
	stopSig bool
	stopped bool
}

func (state *RobotActorV2) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case InitMessageV1:
		state.robot.Init()
	case *actor.Started:
		log.Println("Started, initialize actor here")
	case *actor.Stopping:
		log.Println("Robot is about to be stopped")
	case *actor.Stopped:
		log.Println("Robot Stopped, actor and it's children are stopped")
	case StartMessageV1:
		state.Run(ctx)
	}
}
func (state *RobotActorV2) Init(ctx actor.Context) {

}
func (state *RobotActorV2) markStopped(ctx actor.Context) {
	state.stopped = true
	//ctx.Respond("System Stopped")
}
func (state *RobotActorV2) Stop(ctx actor.Context) {
	state.stopSig = true
}
func (state *RobotActorV2) Run(ctx actor.Context) {
	defer state.markStopped(ctx)
	for {
		if state.stopSig {
			break
		} else {
			//log.Printf("%+v", state.robot.Run())
			state.robot.Run()
		}
	}
}

func NewSystemActorV3() *SystemActorV3 {
	return &SystemActorV3{actor.PIDSet{}, nil}
}

type SystemActorV3 struct {
	actors actor.PIDSet
	w      common.World
}

func (sys *SystemActorV3) Init(ctx actor.Context) {
	sys.w = world.CreateWarehouseWorldWithTaskManager(task.CreateSimulatedTaskManagerSync())
	props := actor.PropsFromProducer(sys.SpawnRobotActor)
	for i := 0; i < 5; i++ {
		sys.actors.Add(ctx.Spawn(props))
	}
	log.Println("System Initialized")
}
func (sys *SystemActorV3) SpawnRobotActor() actor.Actor {
	return &RobotActorV2{robot.NewSimpleWarehouseRobot(uuid.New(), sys.w.GetGraph().Node(1), sys.w), false, false}
}

func (sys *SystemActorV3) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		log.Println("Started, initialize actor here")
		sys.Init(ctx)
		//sys.Run(ctx)
	case StartMessageV1:
		log.Println(ctx.Sender())
		sys.Run(ctx)

	case *actor.Stopping:
		sys.Stop(ctx)
		log.Println("Actor is about to be stopped")
	case *actor.Stopped:
		log.Println("Stopped, actor and it's children are stopped")
	}

	log.Printf("System Received message context %+v", ctx.Message())
}
func (sys *SystemActorV3) Run(ctx actor.Context) {
	log.Println("Running")
	sys.actors.ForEach(
		func(i int, pid actor.PID) {

			ctx.Send(&pid, StartMessageV1{})
		})
	for {
		if sys.w.HasTasks() {

		} else {
			ctx.Respond("Done")
			sys.Stop(ctx)
			break
		}
	}
}
func (sys *SystemActorV3) Stop(ctx actor.Context) {
	sys.actors.ForEach(func(i int, pid actor.PID) {
		ctx.Send(&pid, EndMessageV1{})
		ctx.Stop(&pid)
	})
}
