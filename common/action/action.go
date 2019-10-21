package action

import (
	"maze/common"

	"gonum.org/v1/gonum/graph"
)

type ActionType int

const (
	Pending = iota
	Active
	End
)

type ActionStatus int

const (
	Move = iota
	StartTask
	EndTask
	Act
	Pause
	CustomAction
	NullType
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
	Start  common.Location
	End    common.Location
	Path   []graph.Node
	status ActionStatus
}

func (a *MoveAction) GetChild() Action {
	return a.child
}
func (a *MoveAction) GetType() ActionType {
	return Move
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
func CreateMoveAction(start common.Location, end common.Location) *MoveAction {
	return &MoveAction{nil, start, end, nil, Pending}
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
	here  common.Location
}

func CreateBeginTaskACtion(here common.Location) *BeginTaskAction {
	return &BeginTaskAction{nil, here}
}
func (a *BeginTaskAction) GetChild() Action {
	return a.child
}
func (a *BeginTaskAction) GetType() ActionType {
	return StartTask
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
	here  common.Location
}

func (a EndTaskAction) GetChild() Action {
	return a.child
}
func (a EndTaskAction) GetType() ActionType {
	return EndTask
}
func (a EndTaskAction) HasChild() bool {
	return a.child != nil
}
func (a EndTaskAction) GetContent() interface{} {
	return &a
}
func (a EndTaskAction) SetChild(c Action) {
	a.child = c
}
func CreateEndTaskAction(here common.Location) EndTaskAction {
	return EndTaskAction{nil, here}
}

type NullAction struct {
}

func (n NullAction) GetChild() Action {
	return n
}
func (n NullAction) GetType() ActionType {
	return NullType
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
