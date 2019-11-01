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

func TestV3ActorRun(t *testing.T) {

	ctx := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return impl.NewSystemActorV3()
	})
	pid := ctx.Spawn(props)

	future := ctx.RequestFuture(pid, impl.StartMessageV1{}, 5*time.Second)
	msg, _ := future.Result()

	switch msg.(type) {
	case string:
		if msg != "Done" {
			t.Errorf("Response is %s", msg)
			t.Fail()
		}
	default:
		t.Errorf("Response is %+v", msg)
		t.Fail()
	}
	//
	//select {
	//
	//case <-ctx.Message():
	//
	//case <-time.After(5000 * time.Millisecond):
	//	ctx.Stop(pid)
	//	t.Fail()
	//}
}
