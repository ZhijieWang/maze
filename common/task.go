package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// Task defines the data structure holding the task information
type Task struct {
	id          uuid.UUID
	Origin      graph.Node
	Destination graph.Node
}
