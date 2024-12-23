package api

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/token"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
)

type Server struct{
	config *util.Config
	store *db.Store
	router *gin.Engine
	tokenMaker token.TokenMaker
}


func NewServer(config *util.Config, store *db.Store ) *Server {
	router := gin.Default()
	maker := token.NewPasetoMaker(config.SymmetricKey)

	return &Server{
		config: config,
		store: store,
		router: router,
		tokenMaker: maker,
	}
}

func (server *Server) StartServer (config *util.Config) {
	server.router.Run(config.Address)
}

func (server *Server) SetUpRouter(){
	server.router.POST("/users", server.CreateUser)
	server.router.POST("/users/login", server.loginUser)
	server.router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := server.router.Group("/").Use(authMiddleware(&server.tokenMaker))
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/updateAccout/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccount)

	authRoutes.POST("/transfers", server.createTransfer)

}

func errorResponse(err error) gin.H{
	return gin.H{"err": err.Error()}
}