package main

import (
	"context"
	"database/sql"
	"github.com/314793615/simplebank/api"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/gapi"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {
	config, err := util.NewConfig(".")
	if err != nil {
		log.Fatal("failed to config setting")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connect to db")
	}
	store := db.NewStore(conn)

	server := api.NewServer(config, store)

	server.SetUpRouter()

	server.StartServer(config)

	log.Fatal("server start failed")

}

func runGrpcServer(config util.Config, store *db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create grpc server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("failed to listen on ", config.GrpcServerAddress, err)
	}
	log.Println("grpc server listening on", config.GrpcServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}

}

func runGrpcGateWayServer(config util.Config, store *db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create server: ", err)
	}
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("failed to register grpc server: ", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("failed to listen on ", config.GrpcServerAddress, err)
	}
	err = http.Serve(listener, mux)
	if err != nil {

	}
}
