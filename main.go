package main

import (
	"log"

	"github.com/314793615/simplebank/api"
	"github.com/314793615/simplebank/util"
)


func main(){
	config, err := util.NewConfig(".")
	if err != nil{
		log.Fatal("failed to config setting")
	}

	server := api.NewServer(config)

	server.Set



}