package config

import "os"

var JwtSecret string

func Setup() {
	JwtSecret = os.Getenv("JWT_SECRET")
}
