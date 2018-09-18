package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// Robot is a data holder struct for robot
type Robot struct {
	id       uuid.UUID
	location graph.Node
	task     *Task
	path     []graph.Node
}

// ID returns the robot UUID
func (r *Robot) ID() uuid.UUID {
	return r.id
}
