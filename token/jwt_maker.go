package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}


func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil{
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(maker.secretKey)
	return tokenString, err
}


func (maker *JWTMaker) VerifyToken(tokenString string) (*PayLoad, error) {
	keyFunc := func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &PayLoad{}, keyFunc)
	if err !=  nil {
		return nil, err
	}
	payload, ok := token.Claims.(*PayLoad)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}



