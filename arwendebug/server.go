package arwendebug

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartServer starts a debugging server
func StartServer(facade *DebugFacade, address string) error {
	log.Debug("StartServer()")

	router := gin.Default()

	// TODO: Implement routes (separate PR, when needed)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router.Run(address)
}
