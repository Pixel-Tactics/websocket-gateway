package config

import (
	"os"
)

var RMQUrl string
var RoutingFilePath string
var JwtSecret string
var ParsedRoutes []*Route

func Setup() {
	RoutingFilePath = os.Getenv("ROUTING_FILE_PATH")
	RMQUrl = os.Getenv("RABBITMQ_URL")

	// JwtSecret = os.Getenv("JWT_SECRET")
	JwtSecret = "testos"

	ParsedRoutes = ParseRoutes(RoutingFilePath)
}
