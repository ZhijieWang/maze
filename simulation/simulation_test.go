package simulation_test

import (
	"maze/common"
	"maze/simulation"
	"testing"

	"github.com/google/uuid"
)

type basicObserver struct {
	count int
}

func (b *basicObserver) OnNotify(data interface{}) {

	if data != struct{}{} {
		b.count += 1
	}
}

type traceObserver struct {
	traces []common.Trace
}

func (b *traceObserver) OnNotify(data interface{}) {
	t, ok := data.(common.Trace)
	if ok {
		b.traces = append(b.traces, t)
	}
}

func TestSimulationRunResult(t *testing.T) {
	s := simulation.CreateCentralizedSimulation()
	obs := basicObserver{}
	s.Run(&obs)
	s.Stop()
	if obs.count == 0 {
		t.Errorf("Expect some run. 0 run")
	}

}

func TestSimulationExecuteTask(t *testing.T) {
	s := simulation.CreateCentralizedSimulation()
	obs := traceObserver{}
	s.Run(&obs)
	if len(obs.traces) == 0 {
		t.Error("Failed to capture run traces")
	}
	found := false
	for _, i := range obs.traces {
		t.Logf("Value of trace is %+v ", i)
		if i.RobotID != uuid.Nil {
			found = true
		}
	}
	if !found {
		t.Error("No proper traces were generated from the run")
	}
}
