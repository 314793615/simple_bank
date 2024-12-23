package token

import (
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type PayLoad struct {
	ID uuid.UUID `json:"id"`
	Usename string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayLoad(username string, duration time.Duration) (*PayLoad, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &PayLoad{
		ID: uuid,
		Usename: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (payload *PayLoad) Valid() error {
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}


var _ jwt.Claims = (*PayLoad)(nil)