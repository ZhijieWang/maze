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
	"github.com/google/uuid"
	"log"
	"math/rand"
	"maze/common"
	"maze/common/robot"
	"maze/common/task"
	"maze/common/world"
	"time"
)

type TaskFeederActor struct {
	probability int
	frequency   int
	W           common.World
	comm        chan interface{}
	cap         int
}

func (actor *TaskFeederActor) Run(observer common.Observer) {
	go func() {
		for {
			select {
			default:
				if actor.cap == 0 {
					break
				}
				time.After(5)
				if actor.probability > rand.Intn(100) {
					m := actor.W.GetGraph().Nodes().Len()
					t := task.NewTimePriorityTaskWithParameter(actor.W.GetGraph().Node(int64(rand.Intn(m-1)+1)), actor.W.GetGraph().Node(int64(rand.Intn(m-1)+1)))
					observer.GetChannel() <- fmt.Sprintf("Adding Task %+v", t)

					actor.W.AddTask(t)
					actor.cap--
					log.Printf("Actor Cap is %d", actor.cap)
				}
			case <-actor.comm:
				break
			}

		}
		log.Printf("Break from loop")
	}()
}
func (actor *TaskFeederActor) Init() {

}
func (actor *TaskFeederActor) Stop() {
	close(actor.comm)
}

type ActorRef struct {
	comm  chan interface{}
	robot common.Robot
}

func (actor *ActorRef) Run(observer common.Observer) {

	go func() {
		for {
			select {
			default:
				observer.GetChannel() <- actor.robot.Run()
			case <-actor.comm:
				break
			}
		}
	}()
}
func (actor *ActorRef) Init() {
	actor.robot.Init()

}
func (actor *ActorRef) Stop() {
	close(actor.comm)
}

type System struct {
	W      *world.WarehouseWorld
	stm    *task.SimulatedTaskManagerSync
	refs   []common.Actor
	NumBot int
}

func (s *System) Init() {
	s.stm = task.CreateSimulatedTaskManagerSync()
	s.W = world.CreateWarehouseWorldWithTaskManager(s.stm)

	t := task.NewTimePriorityTask()
	t.Origin = s.W.GetGraph().Node(2)
	t.Destination = s.W.GetGraph().Node(6)

	s.stm.AddTask(t)
	for i := 0; i < s.NumBot; i++ {
		s.refs = append(s.refs, &ActorRef{make(chan interface{}), robot.NewSimpleWarehouseRobot(uuid.New(), s.W.GetGraph().Node(1), s.W)})
	}
	s.refs = append(s.refs, &TaskFeederActor{30, 5, s.W, make(chan interface{}), 5})
	for _, i := range s.refs {
		i.Init()
	}
}
func (s *System) Start(observer common.Observer) {
	for _, i := range s.refs {
		i.Run(observer)
	}
}
func (s System) Stop() bool {
	for _, i := range s.refs {
		i.Stop()
	}
	return true
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
