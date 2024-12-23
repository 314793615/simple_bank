package token

import (
	"time"

	"github.com/o1egl/paseto/v2"
)

type PasetoMaker struct {
	paseto *paseto.V2
	secretKey string
}

func NewPasetoMaker(secretKey string) *PasetoMaker {

	return &PasetoMaker{
		paseto: paseto.NewV2(),
		secretKey: secretKey,
	}
	
}


func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}
	token, err := p.paseto.Encrypt([]byte(p.secretKey), payload, err)
	return token, err
}


func (p *PasetoMaker) VerifyToken(tokenString string) (*PayLoad, error) {
	payload := &PayLoad{}
	err := p.paseto.Decrypt(tokenString, []byte(p.secretKey), payload, nil)
	
	return payload, err
}

var _ TokenMaker = (*PasetoMaker)(nil)