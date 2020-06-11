package arwendebug

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DebugServer is the debugging server
type DebugServer struct {
	facade  *DebugFacade
	address string
}

// NewDebugServer creates a Server object
func NewDebugServer(facade *DebugFacade, address string) *DebugServer {
	return &DebugServer{
		facade:  facade,
		address: address,
	}
}

// StartServer starts the debugging server
func (server *DebugServer) Start() error {
	log.Debug("Start()")

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/account", server.handleCreateAccount)
	router.POST("/deploy", server.handleDeploy)

	// deploy

	// upgrade

	// run

	// query

	return router.Run(server.address)
}

func (server *DebugServer) handleCreateAccount(ginContext *gin.Context) {
	request := CreateAccountRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleCreateAccount.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.CreateAccount(request)
	if err != nil {
		returnBadRequest(ginContext, "handleCreateAccount.CreateAccount", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func (server *DebugServer) handleDeploy(ginContext *gin.Context) {
	request := DeployRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleDeploy.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.DeploySmartContract(request)
	if err != nil {
		returnBadRequest(ginContext, "handleDeploy.DeploySmartContract", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func returnBadRequest(context *gin.Context, errScope string, err error) {
	context.JSON(http.StatusBadRequest, gin.H{
		"error":        fmt.Sprintf("%T", err),
		"errorMessage": err.Error(),
		"errorScope":   errScope,
		"data":         nil,
	})
}

func returnOkResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, gin.H{
		"error":        nil,
		"errorMessage": nil,
		"errorScope":   nil,
		"data":         data,
	})
}
