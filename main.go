package main

import (
	"database/sql"
	"log"

	"github.com/314793615/simplebank/api"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
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