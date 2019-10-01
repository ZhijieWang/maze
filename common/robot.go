package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
)

// RobotID is an alias to UUID for disambiguition purpose
type RobotID = uuid.UUID
type Robot interface {
	ID() RobotID
	Run() Trace
}

// SimpleRobot is a data holder struct for robot
type simpleRobot struct {
	// id is the UUID of the robot
	id RobotID
	// location represents the robot's current location on the graph
	location graph.Node
	// task represents the current work the robot is trying to carry out
	task Task
	// path is the current planned path to deliver the task
	path []graph.Node
}

// ID returns the robot UUID
func (r *simpleRobot) ID() RobotID {
	return r.id
}

// Run is a function to be run by the simulation executor as a go routine
func (r *simpleRobot) Run(w World, tm TaskManager) Trace {
	var tick int
	tick = 1
	if r.task == nil {
		r.task = w.GetTasks()[0]
		w.ClaimTask(r.task.GetTaskID(), r.ID())
		return Trace{
			RobotID:   r.ID(),
			Source:    r.location,
			Target:    r.task.GetDestination(),
			Timestamp: tick,
		}
	}

	if r.location == r.task.GetDestination() {
		// at target location.
		// unset task from robot
		// update task to be done

		//		err := w.TaskUpdate(r.task.GetTaskID(), Completed)
		//		if err != nil {
		//			return Trace{}
		//		}
		return Trace{}
	}
	// go to next location in path
	return Trace{
		RobotID:   r.id,
		Source:    r.location,
		Target:    r.path[0],
		Timestamp: tick,
	}

	//r.localWorld = worldReader.Observe(r.location)

}
func NewSimpleRobot(id RobotID, location graph.Node) Robot {
	s := simpleRobot{
		id,
		location,
		nil,
		nil,
	}
	return &s
}
