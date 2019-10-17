package common

import (
	"log"
	"math/rand"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
)

// TaskGenerator is the generator function for randomly producing tasks
func TaskGenerator(maxTasks int, w World) []Task {

	n := w.GetGraph().Nodes().Len()
	tList := []Task{}
	for i := 0; i < maxTasks; i++ {
		if rand.Intn(2) > 0 {
			uid, _ := uuid.NewUUID()
			tList = append(tList, TimePriorityTask{
				ID:          uid,
				Origin:      w.GetGraph().Node((int64)(rand.Intn(n))),
				Destination: w.GetGraph().Node((int64)(rand.Intn(n))),
			})
		}
	}
	return tList
}

func NoMove(w World, r Robot, t int) Trace {
	return Trace{
		RobotID:   r.ID(),
		Source:    r.Location(),
		Target:    r.Location(),
		Timestamp: t,
	}
}

// // GraphReWeightByRadiation is a graph weight propagation method to recalculate graph edge weight by radiation
// func GraphReWeightByRadiation(world World, trace Trace) {
// 	for _, i := range world.GetRobots() {
// 		world.EdgeWeightPropagation(i.location, 3, 1)
// 	}
// }

// RandMove is a basic function, robot takes a random move that it can move to.
// if there is onlyone path, robot will move
// this is stateless, regardless of previous move taken
func RandMove(w World, r Robot, t int) Trace {
	locs := w.GetGraph().From(r.Location().ID())

	bufs := graph.NodesOf(locs)

	trace := Trace{
		RobotID:   r.ID(),
		Source:    r.Location(),
		Target:    bufs[rand.Intn(len(bufs))],
		Timestamp: t,
	}
	// r.Location() = trace.Target
	return trace
}

// SelectTaskByDistance returns a task from queue, and returns that task. If there is an error, return err
func SelectTaskByDistance(tm PassiveTaskManager, robot Robot, world World) (PriorityTask, []graph.Node, error) {
	tq := tm.GetAllTasks()
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
	pt, _ := path.BellmanFordFrom(robot.Location(), world.GetGraph())
	for _, t := range tq {
		p, weight = pt.To(t.GetOrigination().ID())

		if weight < minWeight {
			minWeight = weight
			tMin = t.(TimePriorityTask)
			pMin = p
		}
	}
	// err := tm.ClaimTask(tMin, robot.ID())
	return tMin, pMin, nil
}
