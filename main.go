package main

import (
	"context"
	"database/sql"
	"github.com/314793615/simplebank/api"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/gapi"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"github.com/314793615/simplebank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"

	"net"
	"net/http"
)

func main() {
	config, err := util.NewConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to config setting")
	}
	log.Info().Str("DRIVER", config.DBDriver).
		Str("DbSource", config.DBSource).Msg("successfully get the config")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}
	store := db.NewStore(conn)
	//redisOpt := asynq.RedisClientOpt{
	//	Addr: config.RedisAddress,
	//}
	//taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	runGinServer(config, store)
	//go runTaskProcessor(redisOpt, store)
	//go runGrpcServer(config, store, taskDistributor)
	//go runGrpcGateWayServer(config, store, taskDistributor)

	log.Info().Msg("server start failed")

}

func runGinServer(config *util.Config, store db.Store) {
	server := api.NewServer(config, store)
	server.SetUpRouter()
	server.StartServer(config)
}

func runGrpcServer(config *util.Config, store db.Store, taskDistributor worker.Distributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen on address %s", config.GrpcServerAddress)
	}
	log.Info().Msgf("grpc server listening on", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start grpc server: ")
	}

}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcess := worker.NewTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcess.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor: ")
	}
}

func runGrpcGateWayServer(config *util.Config, store db.Store, taskDistributor worker.Distributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server: ")
	}
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register grpc server: ")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen on ", config.GrpcServerAddress)
	}
	err = http.Serve(listener, mux)
	if err != nil {

	}
}
