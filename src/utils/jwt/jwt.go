package jwt_auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"pixeltactics.com/websocket-gateway/src/config"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidScheme = errors.New("invalid scheme")
var ErrTokenExpired = errors.New("expired token")

func Validate(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidScheme
		}
		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims.GetSubject()
	}
	return "", ErrInvalidToken
}
