package common

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func TestConcurrentWeightUpdate(t *testing.T) {
	g := NewConcurrentWeightedUndirectedGraph(1.0, 1.0)
	for i := 0; i < 50; i++ {
		go func() {
			g.rwLock.Lock()
			ne := g.NewWeightedEdge(simple.Node(1), simple.Node(2), 0.25)
			g.rwLock.Unlock()
			g.SetWeightedEdge(ne)
		}()
	}

}
