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

package impl_test

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"maze/impl"
	"testing"
	"time"
)

func TestV2ActorRun(t *testing.T) {

	context := actor.EmptyRootContext
	var s = make(chan bool)
	props := actor.PropsFromProducer(func() actor.Actor {
		return impl.NewSystemActorV2(s)
	})
	var pid = context.Spawn(props)
	context.Send(pid, impl.InitMessageV1{})
	context.Send(pid, impl.StartMessageV1{})
	context.Send(pid, impl.EndMessageV1{})

	select {

	case <-s:

	case <-time.After(5000 * time.Millisecond):
		context.Send(pid, impl.ShutdownMessageV1{})
		t.Fail()
	}

}
