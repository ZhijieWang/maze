package common

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
)

// TaskGenerator is the generator function for randomly producing tasks
func TaskGenerator(maxTasks int, w World) []Task {
	n := len(w.GetGraph().Nodes())
	tList := []Task{}
	for i := 0; i < maxTasks; i++ {
		if rand.Intn(2) > 0 {
			uid, _ := uuid.NewUUID()
			tList = append(tList, TimePriorityTask{
				ID:          uid,
				Origin:      w.GetGraph().Nodes()[rand.Intn(n)],
				Destination: w.GetGraph().Nodes()[rand.Intn(n)],
			})
		}
	}
	return tList
}

// GraphReWeightByRadiation is a graph weight propagation method to recalculate graph edge weight by radiation
func GraphReWeightByRadiation(world World, trace Trace) {
	for _, i := range world.GetRobots() {
		world.EdgeWeightPropagation(i.location, 3, 1)
	}
}

// RandMove is a basic function, robot takes a random move that it can move to.
// if there is onlyone path, robot will move
// this is stateless, regardless of previous move taken
func RandMove(w World, r *Robot, t int) Trace {
	locs := w.GetGraph().From(r.location.ID())
	trace := Trace{
		RobotID:   r.ID(),
		Source:    r.location,
		Target:    locs[rand.Intn(len(locs))],
		Timestamp: t,
	}
	r.location = trace.Target
	return trace
}

//TaskMove is a movement policy for Task Oriented movement
func TaskMove(w World, r *Robot, t int) Trace {
	log.Printf("Robot %s can see %d Tasks, current has %vi\n", r.id, len(w.GetTasks()), r.task)
	if r.task != nil {
		log.Printf("Robot %s is carrying out Task %+v\n", r.id, r.task)
		fmt.Printf("%+v\n", r.path)
		fmt.Printf("current location %s, task target location %s\n", r.location, r.task.GetDestination())
		trace := Trace{
			RobotID:   r.ID(),
			Source:    r.location,
			Target:    r.path[0],
			Timestamp: t,
		}
		r.location = r.path[0]
		if len(r.path) == 1 {
			log.Printf("Task %+v done by Robot %s\n", r.task, r.id)
			r.task = nil
			r.path = nil
		} else {
			r.path = r.path[1:]
		}
		return trace
	}
	tasks := w.GetTasks()
	if len(tasks) == 0 {
		log.Println("No Tasks")
		return RandMove(w, r, t)
	}
	tMin := tasks[rand.Intn(len(tasks))]

	pt, _ := path.BellmanFordFrom(r.location, w.GetGraph())
	p, _ := pt.To((*tMin).GetOrigination().ID())
	r.path = p[1:]
	r.task = *tMin
	return Trace{
		RobotID:   r.ID(),
		Source:    r.location,
		Target:    p[0],
		Timestamp: t,
	}
}

// SelectTaskByDistance returns a task from queue, and returns that task. If there is an error, return err
func SelectTaskByDistance(tm PassiveTaskManager, robot Robot, world World) (PriorityTask, []graph.Node, error) {
	tq := tm.GetTasks()
	log.Printf("There are %v tasks currently in Queue\n", len(tq))
	if len(tq) == 0 {
		log.Println("No Tasks")
		return nil, nil, nil
	}
	//claiming task
	var tMin PriorityTask
	minWeight := 0.0
	var weight float64
	var p, pMin []graph.Node
	pt, _ := path.BellmanFordFrom(robot.location, world.GetGraph())
	for _, t := range tq {
		p, weight = pt.To(t.GetOrigination().ID())

		if weight < minWeight {
			minWeight = weight
			tMin = t
			pMin = p
		}
	}
	err := tm.ClaimTask(tMin, robot.ID())
	return tMin, pMin, err
}
