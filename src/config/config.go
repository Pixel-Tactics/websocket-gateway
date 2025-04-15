package config

var JwtSecret string

func Setup() {
	// JwtSecret = os.Getenv("JWT_SECRET")
	JwtSecret = "testos"
}
