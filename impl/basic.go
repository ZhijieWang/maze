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
	"github.com/google/uuid"
	"log"
	"maze/common"
)

type System struct {
	w    *WarehouseWorld
	stm  *SimulatedTaskManager
	refs []*ActorRef
}
type ActorRef struct {
	comm    chan interface{}
	robot   *simpleWarehouseRobot
}

func (actor *ActorRef) Run() {

	go func() {
		for {
			select {
			default:
				actor.robot.Run()
			case <-actor.comm:
				break
			}
		}
		return
	}()
}
func (actor ActorRef) Init() {
	actor.robot.Init()

}
func (s *System) Init() {
	s.w = CreateWarehouseWorld()
	s.stm = CreateSimulatedTaskManager()
	t := common.NewTimePriorityTask()
	t.Origin = s.w.graph.Node(2)
	t.Destination = s.w.graph.Node(6)
	s.stm.AddTask(t)
}
func (s *System) Start() {

	s.refs = append(s.refs, &ActorRef{make(chan interface{}), NewSimpleWarehouseRobot(uuid.New(), s.w.GetGraph().Node(1), s.w, s.stm)})
	for _, i := range s.refs {
		i.Init()
	}
	for _, i := range s.refs {
		go i.Run()
	}
}
func (s System) Stop() {
	for _, i := range s.refs {
		close(i.comm)
	}
}

func (s *System) RunTillStop() {
	for {
		if s.stm.HasTasks() || s.stm.ActiveCount() > 0 {

		} else {
			log.Print("Stopping\n")
			s.Stop()
			break
		}
	}
}
