package simulation

import (
	"log"
	"maze/common"

	"github.com/google/uuid"
)

type Observer interface {
	OnNotify(data interface{})
}
type Event interface {
}
type Notifier interface {
	Register(Observer)
	Deregister(Observer)
	Notify(Event)
}
type Simulation interface {
	Run(obs Observer) error
	Stop() bool
}

type centralizedSimulation struct {
	world common.World
	tm    common.TaskManager
}

func CreateCentralizedSimulation() Simulation {

	var c = centralizedSimulation{}
	c.tm = common.NewBasicTaskManager()
	c.world = common.CreateWorld(1, common.NewBasicTaskManager())
	// c.world = common.CreateBlankWorld()
	var numRobots int = 5
	for i := 0; i < numRobots; i++ {
		rID, err := uuid.NewUUID()
		if err != nil {
			log.Fatal(err)
		}
		c.world.AddRobot(common.NewSimpleRobot(rID, c.world.GetGraph().Nodes().Node()))
	}
	for i := 0; i < 20; i++ {
		t := common.NewTimePriorityTask()
		t.Destination = c.world.GetGraph().Nodes().Node()
		c.tm.AddTask(t)
	}

	return c
}

func (sim centralizedSimulation) Run(obs Observer) error {
	for _, i := range sim.world.GetRobots() {
		trace := i.Run(sim.world, sim.tm)
		sim.world.UpdateRobot(i)
		obs.OnNotify(trace)
	}
	obs.OnNotify(nil)

	return nil
}

func (sim centralizedSimulation) Stop() bool {
	return true
}
