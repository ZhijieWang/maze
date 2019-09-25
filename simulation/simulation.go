package simulation

import (
	"maze/common"
)


type Observer interface{
	OnNotify(data interface{})
}
type Notifier inteface{

	Register(Observer)
	Deregister(Observer)
	Notify(Event)
}
type Simulation interface {
	Run() error
	Run(obs Observer) error
	Stop() bool
}

type centralizedSimulation struct {
	world common.World
}

func CreateCentralizedSimulation() Simulation {

	var c = centralizedSimulation{}
	c.world = common.CreateWorld(1, common.NewBasicTaskManager())
	return c
}

func (sim centralizedSimulation) Run() error {
	return nil

}


func (sim centralizedSimulation) Run(obs Observer) error{
	for _, i := range sim.world.robots {
		i.run()
		obs.notify("Robot with status %d was run", i)
	}
	return nil
}


func (sim centralizedSimulation) Stop() bool {
	return true
}
