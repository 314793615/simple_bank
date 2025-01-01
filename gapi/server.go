package gapi

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/token"
	"github.com/314793615/simplebank/util"
	"github.com/314793615/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          *util.Config
	store           db.Store
	tokenMaker      token.TokenMaker
	redisProcessor  worker.Processor
	taskDistributor worker.Distributor
}

func NewServer(config *util.Config, store db.Store, taskDistributor worker.Distributor) (*Server, error) {
	tokenMaker := token.NewPasetoMaker(config.TokenSymmetricKey)
	return &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
		//RedisProcessor: worker.NewTaskProcessor(config.RedisAddress),
	}, nil
}
