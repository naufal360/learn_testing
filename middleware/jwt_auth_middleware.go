package middleware

import (
	"learn_testing/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userId int, name string) (string, error) {
	key := []byte(config.ViperEnvVariable("SECRET_KEY"))
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
