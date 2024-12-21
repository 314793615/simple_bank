package api

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
)

type Server struct{
	config *util.Config
	store *db.Store
	router *gin.Engine
}


func NewServer(config *util.Config ) *Server {
	d := sql.Open(config.DBDriver, config.DBSource)
	store := db.NewStore(d)
	router := gin.New()
	return &Server{
		config: config,
		store: store,
		router: router,
	}
}

func (server *Server) StartServer (config util.Config) {
	server.router.Run(config.Address)
}

func (server *Server) SetUpRouter(){
	server.router.POST("/createAccount", server.CreateAccount)
	server.router.POST("/updateAccout", server.UpdateAccount)
}

func errorResponse(err error) gin.H{
	return gin.H{"err": err.Error()}
}