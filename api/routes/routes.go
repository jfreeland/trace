package routes

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Fallthrough provides a route back to index.html
func Fallthrough(rootDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") || path.Ext(c.Request.URL.Path) != "" {
			c.JSON(404, gin.H{"message": "does not exist"})
			c.Abort()
		} else {
			http.ServeFile(c.Writer, c.Request, fmt.Sprintf("%s/index.html", rootDir))
			c.Abort()
		}
	}
}
