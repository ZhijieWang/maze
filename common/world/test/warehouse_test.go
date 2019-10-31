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

package test

import (
	"maze/common/world"
	"testing"
)

func TestBeAbleToGetRobots(t *testing.T) {
	w := world.CreateWarehouseWorld()

	for _, r := range w.GetRobots() {
		if r == nil {
			t.Errorf("We got a problem, robot shouldn't be nil")
			t.FailNow()
		}
	}
}
