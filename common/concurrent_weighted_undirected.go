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

// SetWeightedEdge mimics the original implementation with mutex
func (g *ConcurrentWeightedUndirectedGraph) SetWeightedEdge(e graph.WeightedEdge) {
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

// NewWeightedEdge returns a new weighted edge from the source to the destination node.
func (g *ConcurrentWeightedUndirectedGraph) NewWeightedEdge(from, to graph.Node, weight float64) graph.WeightedEdge {
	g.rwLock.Lock()
	defer g.rwLock.Unlock()
	return g.WeightedUndirectedGraph.NewWeightedEdge(from, to, weight)
}
