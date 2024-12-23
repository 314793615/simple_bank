package token

import (
	"time"
)

type TokenMaker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(tokenString string) (*PayLoad, error)
}