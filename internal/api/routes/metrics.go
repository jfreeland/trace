package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jfreeland/trace/internal/storage"
)

type Metric struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type MetricsResult struct {
	Nodes []*Metric
}

// GetMetrics returns metrics for a TraceHost
func GetMetrics(db storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("TraceHost")
		results, ok := db.GetResults(host)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"status": "data not in store"})

			return
		}
		var metrics []*Metric
		for _, result := range results {
			lastHop := result.Hops[len(result.Hops)-1]
			metrics = append(metrics, &Metric{
				Date:  result.Time.UTC().Format("2006-01-02T15:04:05.999Z"),
				Value: lastHop.Duration.Milliseconds(),
			})
		}
		c.JSON(http.StatusOK, metrics)
	}
}
