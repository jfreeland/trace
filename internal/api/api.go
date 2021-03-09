package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/jfreeland/trace/internal/api/routes"
	"github.com/jfreeland/trace/internal/storage"
	"github.com/jfreeland/trace/internal/tracer"
)

// New creates a gin engine to serve HTTP requests.
func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{http.MethodPost, http.MethodPut, http.MethodGet},
		}),
	)

	return router
}

// AddRoutes adds routes to our router.
func AddRoutes(router *gin.Engine, db storage.Storage, tracer tracer.Tracer) {
	apiRoutes := router.Group("/v0")
	// routes.GET("/", ReturnData(db))
	// routes.GET("/key/:key", ReturnSingleObject(db))
	apiRoutes.GET("/graph", routes.GetGraph(db))
	apiRoutes.GET("/metrics", routes.GetMetrics(db))
	apiRoutes.POST("/start", routes.AddTrace(db, tracer))
	apiRoutes.DELETE("/stop", routes.StopTrace(db, tracer))
}

func AddStatic(router *gin.Engine) {
	const staticDir = "/ui"
	router.Use(static.Serve("/", static.LocalFile(staticDir, true)))
	router.Use(routes.Fallthrough(staticDir))
}

// ReturnData handles HTTP GETs.
// func ReturnData(db storage.Storage, tracer tracer.Tracer) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.JSON(http.StatusOK, db.GetAll())
// 	}
// }

// AddData handles HTTP PUTs and POSTs.
//
// TODO: This is a very crude implementation.  I could (should?) look at the
// Content-Type header and decode based on the input type.  For the sake of
// expediency I've chosen to only support JSON input and if it's not JSON I'll
// take the input as a string and assign you a key of now.
// func AddData(db storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer c.Request.Body.Close()
// 		input, err := ioutil.ReadAll(c.Request.Body)
// 		v := string(input)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not read body"})

// 			return
// 		}
// 		data := make(map[string]interface{})
// 		if err := json.NewDecoder(strings.NewReader(v)).Decode(&data); err != nil {
// 			log.Debugf("data was not a json object: %v", err)
// 			data[time.Now().UTC().String()] = v
// 			db.Store(&data)
// 			c.JSON(http.StatusOK, gin.H{"status": "data was not json, key is current timestamp"})

// 			return
// 		}
// 		db.Store(&data)
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	}
// }

// ReturnSingleObject handles returning a single object from the map.
// func ReturnSingleObject(db storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		key := c.Param("key")
// 		value, ok := db.Get(key)
// 		if !ok {
// 			c.JSON(http.StatusOK, gin.H{"status": "data not in store"})

// 			return
// 		}
// 		c.String(http.StatusOK, "%v", value)
// 	}
// }
