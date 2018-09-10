package participants

import (
	"gonum.org/v1/gonum/graph"
	"log"
	"net"
	"net/rpc"
)

func init() {

	server := rpc.NewServer()
	g := concrete.NewGonumGraph(true)
	var n0, n1, n2, n3 concrete.GonumNode = 0, 1, 2, 3
	g.AddNode(n0, []graph.Node{n1, n2})
	g.AddEdge(concrete.GonumEdge{n2, n3})
	path, v := search.BreadthFirstSearch(n0, n3, g)
	fmt.Println("path:", path)
	fmt.Println("nodes visited:", v)
	// registerArith(server, arith)

	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	server.Accept(l)
}
