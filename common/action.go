package common

import (
	"gonum.org/v1/gonum/graph"
)

type ActionType int

const (
	PendingStatus = iota
	ActiveStatus
	EndStatus
)

type ActionStatus int

const (
	ActionTypeMove = iota
	ActionTypeStartTask
	ActionTypeEndTask
	ActionTypeNull
)

type Action interface {
	GetType() ActionType
	GetContent() interface{}
	HasChild() bool
	GetChild() Action
	SetChild(Action)
}

type DurationAction interface {
	GetStatus() ActionStatus
	SetStatus(ActionStatus)
}

type MoveAction struct {
	child  Action
	Start  Location
	End    Location
	Path   []graph.Node
	status ActionStatus
}

func (a *MoveAction) GetChild() Action {
	return a.child
}
func (a *MoveAction) GetType() ActionType {
	return ActionTypeMove
}
func (a *MoveAction) HasChild() bool {
	return a.child != nil
}
func (a *MoveAction) SetChild(c Action) {
	a.child = c
}
func (a *MoveAction) GetContent() interface{} {
	return a
}
func CreateMoveAction(start Location, end Location) *MoveAction {
	return &MoveAction{nil, start, end, nil, PendingStatus}
}
func CreateMoveActionWithPath(start, end Location, path []graph.Node) *MoveAction {
	return &MoveAction{nil, start, end, path, PendingStatus}
}
func (a *MoveAction) GetStatus() ActionStatus {
	return a.status
}
func (a *MoveAction) SetStatus(s ActionStatus) {
	a.status = s
}

// func Plan() Action {
// 	return MoveAction{nil, nil, nil}
// }

type BeginTaskAction struct {
	child Action
	here  Location
}

func CreateBeginTaskAction(here Location) *BeginTaskAction {
	return &BeginTaskAction{nil, here}
}
func (a *BeginTaskAction) GetChild() Action {
	return a.child
}
func (a *BeginTaskAction) GetType() ActionType {
	return ActionTypeStartTask
}
func (a *BeginTaskAction) HasChild() bool {
	return a.child != nil
}
func (a *BeginTaskAction) GetContent() interface{} {
	return &a
}
func (a *BeginTaskAction) SetChild(c Action) {
	a.child = c
}

type EndTaskAction struct {
	child Action
	here  Location
}

func (a *EndTaskAction) GetChild() Action {
	return a.child
}
func (a *EndTaskAction) GetType() ActionType {
	return ActionTypeEndTask
}
func (a *EndTaskAction) HasChild() bool {
	return a.child != nil
}
func (a *EndTaskAction) GetContent() interface{} {
	return &a
}
func (a *EndTaskAction) SetChild(c Action) {
	a.child = c
}
func CreateEndTaskAction(here Location) *EndTaskAction {
	return &EndTaskAction{nil, here}
}

type NullAction struct {
}

func (n NullAction) GetChild() Action {
	return n
}
func (n NullAction) GetType() ActionType {
	return ActionTypeNull
}

func (n NullAction) HasChild() bool {
	return false
}

func (n NullAction) GetContent() interface{} {
	return n
}

func (n NullAction) SetChild(a Action) {

}

func Null() NullAction {
	return NullAction{}
}
