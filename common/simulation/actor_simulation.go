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
	"fmt"
	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/google/uuid"
	"log"
	"maze/common"
	"maze/common/robot"
	"maze/common/world"
)

type Hello struct{ Who string }
type HelloActor struct{}

func (state *HelloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

type RobotActor struct {
	robot common.Robot
}

func (state *RobotActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case InitMessage:
		state.robot.Init()

	case StartMessage:
		log.Printf("%+v", state.robot.Run())
	}
}

type InitMessage struct {
}
type StartMessage struct {
}

func Run() {
	//context := actor.EmptyRootContext
	//props := actor.PropsFromProducer(func() actor.Actor { return &HelloActor{} })
	//var pid = context.Spawn(props)
	//context.Send(pid, Hello{Who: "Roger"})
	//console.ReadLine()
	context := actor.EmptyRootContext
	var w common.World
	w = world.CreateWarehouseWorld()
	props := actor.PropsFromProducer(func() actor.Actor {

		rid, _ := uuid.NewUUID()
		return &RobotActor{robot.NewSimpleWarehouseRobot(rid, w.GetGraph().Node(1), w)}
	})
	var pid = context.Spawn(props)
	context.Send(pid, InitMessage{})
	context.Send(pid, StartMessage{})
	console.ReadLine()

}
