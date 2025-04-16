package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const (
	ROUTE_QUEUE        = "queue"
	ROUTE_STREAM       = "stream"
	DIRECTION_INCOMING = "incoming"
	DIRECTION_OUTGOING = "outgoing"
)

type RouteType = string
type Direction = string

type Routes struct {
	Routes []*Route `json:"routes" validate:"required,min=1,dive"`
}

type Route struct {
	Name      string    `yaml:"name" validate:"required,min=1"`
	Prefix    string    `yaml:"prefix" validate:"required,min=1"`
	Type      RouteType `yaml:"type" validate:"required,oneof=queue stream"`
	Direction Direction `yaml:"direction" validate:"required,oneof=incoming outgoing"`
	Schema    string    `yaml:"schema"`
}

func ParseRoutes(filepath string) []*Route {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var routes Routes
	err = yaml.Unmarshal(data, &routes)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(routes); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("validation error: Field %s failed validation (%s)\n",
				err.Field(), err.Tag())
		}
		log.Fatalln("terminated")
	}
	return routes.Routes
}
