package config

const (
	ROUTE_QUEUE  = "ROUTE_QUEUE"
	ROUTE_STREAM = "ROUTE_STREAM"
)

type RouteType = string

type Route struct {
	Name   string
	Prefix string
	Type   RouteType
}

var Routes = []Route{
	{
		Name:   "session-service",
		Prefix: "session/",
		Type:   ROUTE_QUEUE,
	},
}
