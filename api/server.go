package api

import (
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/token"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     *util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
}

func NewServer(config *util.Config, store db.Store) *Server {
	router := gin.Default()
	maker := token.NewPasetoMaker(config.TokenSymmetricKey)

	return &Server{
		config:     config,
		store:      store,
		router:     router,
		tokenMaker: maker,
	}
}

func (server *Server) StartServer(config *util.Config) {
	server.router.Run(config.GinServerAddress)
}

func (server *Server) SetUpRouter() {
	server.router.POST("/users", server.CreateUser)
	server.router.POST("/users/login", server.LoginUser)
	server.router.POST("/tokens/renew_access", server.RenewAccessToken)
	server.router.GET("/accounts/:id", server.GetAccount)

	//authRoutes := server.router.Group("/").Use(authMiddleware(server.tokenMaker))
	//authRoutes.POST("/accounts", server.CreateAccount)
	//authRoutes.GET("/accounts/:id", server.GetAccount)
	//authRoutes.GET("/accounts", server.ListAccount)

}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
