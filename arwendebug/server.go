package arwendebug

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartServer -
func StartServer(facade *DebugFacade, address string) error {
	log.Debug("StartServer()")

	// TODO: gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router.Run(address)
}
