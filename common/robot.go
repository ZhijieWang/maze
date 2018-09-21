package common

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// RobotID is an alias to UUID for disambiguition purpose
type RobotID = uuid.UUID

// Robot is a data holder struct for robot
type Robot struct {
	// id is the UUID of the robot
	id RobotID
	// location represents the robot's current location on the graph
	location graph.Node
	// task represents the current work the robot is trying to carry out
	task Task
	// path is the current planned path to deliver the task
	path       []graph.Node
	localWorld simple.WeightedUndirectedGraph
}

// ID returns the robot UUID
func (r *Robot) ID() RobotID {
	return r.id
}

// Run is a function
func (r Robot) Run(clock chan int, policy func(), worldReader World, taskReader PassiveTaskManager) Trace {
	var tick int
	tick = <-clock
	if r.task == nil {
		task, p, _ := SelectTaskByDistance(taskReader, r, worldReader)
		r.task = task
		return Trace{
			RobotID:   r.ID(),
			Source:    r.location,
			Target:    p[0],
			Timestamp: tick,
		}
	}

	if r.location == r.task.GetDestination() {
		// at target location.
		// unset task from robot
		// update task to be done

		err := taskReader.TaskUpdate(r.task.GetTaskID(), Completed)
		if err != nil {
			return Trace{}
		}
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
