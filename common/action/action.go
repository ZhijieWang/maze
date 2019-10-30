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

package action

import (
	"gonum.org/v1/gonum/graph"
	"maze/common"
	"reflect"
)

type MoveAction struct {
	child  common.Action
	Start  common.Location
	End    common.Location
	Path   []graph.Node
	status common.ActionStatus
}

func (a *MoveAction) GetChild() common.Action {
	return a.child
}
func (a *MoveAction) GetType() common.ActionType {
	return common.ActionTypeMove
}
func (a *MoveAction) HasChild() bool {
	return a.child != nil
}
func (a *MoveAction) SetChild(c common.Action) {
	a.child = c
}
func (a *MoveAction) GetContent() interface{} {
	return a
}
func (a *MoveAction) Equal(other common.Action) bool {
	if a.GetType() == a.GetType() {
		cast := other.(*MoveAction)
		return cast.Start == a.Start && cast.End == a.End && reflect.DeepEqual(cast.Path, a.Path) && cast.child.Equal(a.child)
	} else {
		return false
	}
}

func CreateMoveAction(start common.Location, end common.Location) *MoveAction {
	return &MoveAction{nil, start, end, nil, common.PendingStatus}
}
func CreateMoveActionWithPath(start, end common.Location, path []graph.Node) *MoveAction {
	return &MoveAction{nil, start, end, path, common.PendingStatus}
}
func (a *MoveAction) GetStatus() common.ActionStatus {
	return a.status
}
func (a *MoveAction) SetStatus(s common.ActionStatus) {
	a.status = s
}

type BeginTaskAction struct {
	child common.Action
	here  common.Location
}

func CreateBeginTaskAction(here common.Location) *BeginTaskAction {
	return &BeginTaskAction{nil, here}
}
func (a *BeginTaskAction) GetChild() common.Action {
	return a.child
}
func (a *BeginTaskAction) GetType() common.ActionType {
	return common.ActionTypeStartTask
}
func (a *BeginTaskAction) HasChild() bool {
	return a.child != nil
}
func (a *BeginTaskAction) GetContent() interface{} {
	return &a
}
func (a *BeginTaskAction) SetChild(c common.Action) {
	a.child = c
}
func (a *BeginTaskAction) Equal(other common.Action) bool {
	if a.GetType() == a.GetType() {
		cast := other.(*BeginTaskAction)
		return cast.here == a.here && cast.child.Equal(a.child)
	} else {
		return false
	}
}

type EndTaskAction struct {
	child common.Action
	here  common.Location
}

func (a *EndTaskAction) GetChild() common.Action {
	return a.child
}
func (a *EndTaskAction) GetType() common.ActionType {
	return common.ActionTypeEndTask
}
func (a *EndTaskAction) HasChild() bool {
	return a.child != nil
}
func (a *EndTaskAction) GetContent() interface{} {
	return &a
}
func (a *EndTaskAction) SetChild(c common.Action) {
	a.child = c
}
func CreateEndTaskAction(here common.Location) *EndTaskAction {
	return &EndTaskAction{nil, here}
}
func (a *EndTaskAction) Equal(other common.Action) bool {
	if a.GetType() == a.GetType() {
		cast := other.(*EndTaskAction)
		return cast.here == a.here && cast.child.Equal(a.child)
	} else {
		return false
	}
}

type NullAction struct {
}

func (n NullAction) GetChild() common.Action {
	return n
}
func (n NullAction) GetType() common.ActionType {
	return common.ActionTypeNull
}

func (n NullAction) HasChild() bool {
	return false
}

func (n NullAction) GetContent() interface{} {
	return n
}

func (n NullAction) SetChild(a common.Action) {

}

func Null() NullAction {
	return NullAction{}
}
func (n NullAction) Equal(other common.Action) bool {
	return n.GetType() == n.GetType()

}
