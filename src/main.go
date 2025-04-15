package main

import (
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/gateway"
	"pixeltactics.com/websocket-gateway/src/websockets"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Setup()

	gatewayRouter := gateway.NewRouter()
	clientHub := websockets.NewClientHub(gatewayRouter)

	go clientHub.Run()

	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		websockets.ServeWebSocket(clientHub, context.Writer, context.Request)
	})

	router.Run("0.0.0.0:8080")
}
