package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jfreeland/trace/internal/storage"
	"github.com/jfreeland/trace/internal/tracer"
)

// AddTrace adds a new traceroute
func AddTrace(db storage.Storage, tracer tracer.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("TraceHost")
		go tracer.Run(host)
	}
}

// StopTrace stops a running traceroute
func StopTrace(db storage.Storage, tracer tracer.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("TraceHost")
		tracer.Stop(host)
	}
}
