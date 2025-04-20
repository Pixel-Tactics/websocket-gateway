package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/events"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
	"pixeltactics.com/websocket-gateway/src/router"
	"pixeltactics.com/websocket-gateway/src/websockets"
)

func main() {
	godotenv.Load()
	config.Setup()

	rmqManager := communication.NewRMQManager()

	eventManager := events.NewAsyncEventManager()
	controlRouter := router.NewControlRouter()
	clientHub := websockets.NewClientHub(controlRouter, eventManager)
	inFactory := router.NewIncomingRouterFactory(rmqManager)
	outFactory := router.NewOutgoingRouterFactory(rmqManager, eventManager, clientHub)
	for _, routeConfig := range config.ParsedRoutes {
		if routeConfig.Direction == config.DIRECTION_INCOMING {
			router, err := inFactory.Generate(routeConfig)
			if err != nil {
				panic(err)
			}
			controlRouter.AddIncomingRouter(router)
		} else {
			router, err := outFactory.Generate(routeConfig)
			if err != nil {
				panic(err)
			}
			controlRouter.AddOutgoingRouter(router)
		}
	}

	go eventManager.Run()
	controlRouter.Run()
	go clientHub.Run()

	ginRouter := gin.Default()
	ginRouter.GET("/", func(context *gin.Context) {
		websockets.ServeWebSocket(clientHub, context.Writer, context.Request)
	})
	ginRouter.Run("0.0.0.0:8080")
}
