package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jfreeland/trace/internal/storage"
)

type GraphNode struct {
	ID string `json:"id"`
}

type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	// TODO: Add weight.
}

type GraphResult struct {
	Nodes []*GraphNode `json:"nodes"`
	Links []*GraphEdge `json:"links"`
}

// GetGraph returns a graph of a traceroute
func GetGraph(db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("TraceHost")
		results, ok := db.GetResults(host)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"status": "data not in store"})

			return
		}
		var nodes []*GraphNode
		var edges []*GraphEdge
		for _, result := range results {
			hopCount := len(result.Hops)
			for idx, hop := range result.Hops {
				if !nodeInNodes(hop.Host.IP, nodes) {
					node := &GraphNode{
						ID: hop.Host.IP,
					}
					nodes = append(nodes, node)
				}
				if idx < hopCount-1 {
					edge := &GraphEdge{
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

func nodeInNodes(n string, nodes []*GraphNode) bool {
	for _, node := range nodes {
		if n == node.ID {
			return true
		}
	}
	return false
}

func edgeInEdges(e *GraphEdge, edges []*GraphEdge) bool {
	for _, edge := range edges {
		if (e.Source == edge.Source && e.Target == edge.Target) || (e.Source == edge.Target && e.Target == edge.Source) {
			return true
		}
	}
	return false
}
