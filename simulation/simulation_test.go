package simulation_test

import (
	"maze/simulation"
	"testing"
)

func TestSimulationRunStart(t *testing.T) {
	var s = simulation.CreateCentralizedSimulation()
	err := s.Run()
	if err != nil {
		t.Error("The run failed to start")
	}
}

type baseObserver struct {
	count  int
	Output chan int
}

func (b *baseObserver) OnNotify(data interface{}) {
	if data != struct{}{} {
		b.count += 1
	} else {
		b.Output <- b.count
	}
}
func TestSimulationRunResult(t *testing.T) {
	var s = simulation.CreateCentralizedSimulation()
	c := make(chan int)
	count := <-c

	obs := baseObserver{0, c}
	s.Run(&obs)

	if count == 0 {
		t.Errorf("Expect some run. 0 run")
	}

}
