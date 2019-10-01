package common

import (
	"gonum.org/v1/gonum/graph"
)

// Trace is data structure to hold data that can be used for path planning
type Trace struct {
	RobotID   RobotID
	Source    graph.Node
	Target    graph.Node
	Timestamp int
}
