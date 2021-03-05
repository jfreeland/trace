package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jfreeland/trace/internal/storage"
)

type Node struct {
	ID string `json:"id"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	// TODO: Add weight.
}

type GraphResult struct {
	Nodes []*Node `json:"nodes"`
	Links []*Edge `json:"links"`
}

// GetGraph adds a new traceroute
func GetGraph(db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("TraceHost")
		results, ok := db.GetResults(host)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"status": "data not in store"})

			return
		}
		var nodes []*Node
		var edges []*Edge
		for _, result := range results {
			hopCount := len(result.Hops)
			for idx, hop := range result.Hops {
				if !nodeInNodes(hop.Host.IP, nodes) {
					node := &Node{
						ID: hop.Host.IP,
					}
					nodes = append(nodes, node)
				}
				if idx < hopCount-1 {
					edge := &Edge{
						Source: hop.Host.IP,
						Target: result.Hops[idx+1].Host.IP,
					}
					if !edgeInEdges(edge, edges) {
						edges = append(edges, edge)
					}
				}
			}
		}
		c.JSON(http.StatusOK, &GraphResult{
			Nodes: nodes,
			Links: edges,
		})
	}
}

func nodeInNodes(n string, nodes []*Node) bool {
	for _, node := range nodes {
		if n == node.ID {
			return true
		}
	}
	return false
}

func edgeInEdges(e *Edge, edges []*Edge) bool {
	for _, edge := range edges {
		if (e.Source == edge.Source && e.Target == edge.Target) || (e.Source == edge.Target && e.Target == edge.Source) {
			return true
		}
	}
	return false
}
