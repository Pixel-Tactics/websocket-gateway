package main

import (
	"github.com/joho/godotenv"
	"pixeltactics.com/websocket-gateway/src/config"
)

func main() {
	godotenv.Load()
	config.Setup()

	// gatewayRouter := gateway.NewRouter()
	// clientHub := websockets.NewClientHub(gatewayRouter)

	// go clientHub.Run()

	// router := gin.Default()

	// router.GET("/", func(context *gin.Context) {
	// 	websockets.ServeWebSocket(clientHub, context.Writer, context.Request)
	// })

	// router.Run("0.0.0.0:8080")
}
