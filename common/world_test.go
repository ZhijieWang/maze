package common

import (
	"testing"
)

func TestConcurrentWorldSimulation(t *testing.T) {
	w := CreateWorld(2, true)
	for i := 0; i < 10000; i++ {
		go func() {
			w.UpdateWeight(w.grid.WeightedEdgeBetween(1, 2), 0.5)
		}()
	}
}
