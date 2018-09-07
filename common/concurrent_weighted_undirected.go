package common

import (
	"sync"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// ConcurrentWeightedUndirectedGraph extende the GoNum WeightedUndirectedGraph with sync.Mutex for concurrent access
type ConcurrentWeightedUndirectedGraph struct {
	*simple.WeightedUndirectedGraph
	rwLock *sync.RWMutex
}

// ConcurrentSetWeightedEdge mimics the original implementation with mutex
func (g *ConcurrentWeightedUndirectedGraph) ConcurrentSetWeightedEdge(e graph.WeightedEdge) {
	g.rwLock.Lock()
	g.WeightedUndirectedGraph.SetWeightedEdge(e)
	g.rwLock.Unlock()
}

// NewConcurrentWeightedUndirectedGraph returns an WeightedUndirectedGraph with the specified self and absent
// edge weight values.
func NewConcurrentWeightedUndirectedGraph(self, absent float64) *ConcurrentWeightedUndirectedGraph {
	return &ConcurrentWeightedUndirectedGraph{
		simple.NewWeightedUndirectedGraph(self, absent),
		&sync.RWMutex{},
	}
}

// WeightedEdgeBetween returns the weighted edge between nodes x and y.
func (g *ConcurrentWeightedUndirectedGraph) WeightedEdgeBetween(xid, yid int64) graph.WeightedEdge {
	g.rwLock.Lock()
	defer g.rwLock.Unlock()
	return g.WeightedUndirectedGraph.WeightedEdgeBetween(xid, yid)
}
