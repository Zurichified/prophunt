package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"os/exec"
	"strings"
)

// We need to give every funccall in CG an int representation to later build edges
type Node graph.Node

// From the funccall to the according Node
type NodeMap map[string]Node

var DirectedGraph = simple.NewDirectedGraph()

type ReachableNodes map[Node]int

// Function that executes callgraph tool on the project path and gives the edges in form of STDOUT
func CallgraphCommand(path string) ([]byte, error) {
	//args := "-algo pta -test -format digraph"
	gopath := strings.Split(path, string(os.PathSeparator)+"src")[0]
	os.Setenv("GOPATH", gopath)

	edges := exec.Command("callgraph", "-algo", "pta", "-test", "-format", "digraph", path)
	edges.Env = os.Environ()

	return edges.Output()
}

// Function that returns a directed Graph to look for reachable nodes from a node
func GetCallgraph(result []byte, callgraph *simple.DirectedGraph, nodeMap NodeMap) (*simple.DirectedGraph, NodeMap) {

	stdout := string(result)
	lines := strings.Split(stdout, "\n")

	for _, edge := range lines {
		// If cmd gave an empty string
		if len(edge) == 0 {
			continue
		}

		funcs := strings.Split(edge, " ")
		caller := funcs[0]
		callee := funcs[1]

		var callerNode, calleeNode Node

		// Setting caller node
		if nodeMap[caller] != nil {
			callerNode = nodeMap[caller]
		} else {
			callerNode = callgraph.NewNode()
			callgraph.AddNode(callerNode)
			nodeMap[caller] = callerNode
		}

		// Setting callee node
		if nodeMap[callee] != nil {
			calleeNode = nodeMap[callee]
		} else {
			calleeNode = callgraph.NewNode()
			callgraph.AddNode(calleeNode)
			nodeMap[callee] = calleeNode
		}

		// Graph does not allow and edge to the same Node, hence, skipping if that's the case
		if callerNode == calleeNode {
			continue
		}

		// Creating an edge between caller and callee and feeding this into the graph
		newedge := callgraph.NewEdge(callerNode, calleeNode)
		callgraph.SetEdge(newedge)
	}
	return callgraph, nodeMap
}

// Function that gives all the reachable nodes from a Node
func FromNext(cg *simple.DirectedGraph, node Node, reachables ReachableNodes) {
	// Stop if node is visited already
	if reachables[node] == 1 {
		return
	} else { reachables[node] = 1
		//fmt.Println("printing: " + strconv.FormatInt(node.ID(), 10))
	}

	// Stop if there are no other funccalls from given node
	result := cg.From(node.ID())
	if result.Len() == 0 {
		return
	}

	// Main recursion loop
	for {
		if result.Next() != false {
			var nextNode Node
			nextNode = result.Node()
			FromNext(cg, nextNode, reachables)
		} else {
			break
		}
	}
}

func NodeLookUp(node Node, nodemap NodeMap) string {
	for key, value := range nodemap {
		if node == value {
			return key
		}
	}
	return ""
}


