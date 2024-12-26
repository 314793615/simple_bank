package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/314793615/simplebank/api"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/gapi"
	"github.com/314793615/simplebank/pb"
	"github.com/314793615/simplebank/util"
	"google.golang.org/grpc"
)


func main(){
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


















































































func runGrpcServer(config *util.Config, store *db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gapi server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	listen, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("cannot create grpc listen: ", err)
	}

	http.Serve(listen, grpcServer)

}

func runGrpcGateWayServer(config *util.Config, store *db.Store) {
	server, err := gapi.NewServer()
}