package simulation_test

import (
	"maze/simulation"
	"testing"
)

type basicObserver struct {
	count int
}

func (b *basicObserver) OnNotify(data interface{}) {

	if data != struct{}{} {
		b.count += 1
	}
}
func TestSimulationRunResult(t *testing.T) {
	s := simulation.CreateCentralizedSimulation()
	obs := basicObserver{}
	s.Run(&obs)

	if obs.count == 0 {
		t.Errorf("Expect some run. 0 run")
	}

}
