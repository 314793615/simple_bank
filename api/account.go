package api

import (
	"net/http"
	"time"

	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)


type createAccountParams struct {
	Owner string `json:"owner" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) CreateAccount(ctx *gin.Context){
	var arg createAccountParams
	var err error
	err = ctx.ShouldBindJSON(&arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner: arg.Owner,
		Balance: arg.Balance,
		Currency: arg.Currency,
	})
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return 
	}

	ctx.JSON(http.StatusOK, account)

} 

func (Server *Server) UpdateAccount(ctx *gin.Context){


}

