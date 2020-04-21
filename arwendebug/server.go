package arwendebug

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func startServer(address string) {
	// TODO: gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Run(address)
}
