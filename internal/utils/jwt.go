package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

var jwtKey = os.Getenv("JWT_KEY")

type claims struct {
	*model.UserPayload
	jwt.StandardClaims
}

func CreateToken(userPayload *model.UserPayload) (string, error) {
	return CreateTokenWithExpire(userPayload, time.Now().Add(336*time.Hour).Unix())
}

func CreateTokenWithExpire(userPayload *model.UserPayload, exp int64) (string, error) {
	sub := strconv.Itoa(userPayload.ID)
	claims := &claims{
		UserPayload: userPayload,
		StandardClaims: jwt.StandardClaims{
			Subject:   sub,
			ExpiresAt: exp,
			IssuedAt:  time.Now().Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at.Header["typ"] = "JWT"
	token, err := at.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(token string) (*jwt.Token, *model.UserPayload, error) {
	claims := &claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	return parsedToken, claims.UserPayload, err
}
