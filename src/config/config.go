package config

import (
	"os"
)

var RoutingPath string
var JwtSecret string
var ParsedRoutes []*Route

func Setup() {
	RoutingPath = os.Getenv("ROUTING_PATH")

	// JwtSecret = os.Getenv("JWT_SECRET")
	JwtSecret = "testos"

	ParsedRoutes = ParseRoutes(RoutingPath)
}
