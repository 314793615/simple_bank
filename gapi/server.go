package gapi

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/token"
	"github.com/314793615/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config *util.Config
	store *db.Store
	tokenMaker token.TokenMaker
}


func NewServer(config *util.Config, store *db.Store) (*Server, error) {
	tokenMaker:= token.NewPasetoMaker(config.TokenSymmetricKey)
	return &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}, nil

}
