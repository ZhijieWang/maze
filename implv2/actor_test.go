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

package implv2_test

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"maze/implv2"
	"testing"
	"time"
)

func TestActorRun(t *testing.T) {

	context := actor.EmptyRootContext
	var s = make(chan bool)
	props := actor.PropsFromProducer(func() actor.Actor {
		return implv2.NewSystemActor(s)
	})
	var pid = context.Spawn(props)
	context.Send(pid, implv2.InitMessage{})
	context.Send(pid, implv2.StartMessage{})
	context.Send(pid, implv2.EndMessage{})

	select {

	case <-s:

	case <-time.After(5000 * time.Millisecond):
		context.Send(pid, implv2.ShutdownMessage{})
		t.Fail()
	}

}
