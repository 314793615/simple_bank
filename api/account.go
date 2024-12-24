package api

import (
	"database/sql"
	"net/http"
	
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner: arg.Owner,
		Balance: arg.Balance,
		Currency: arg.Currency,
	})
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, account)

} 

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`

}

func (server *Server) GetAccount(ctx *gin.Context){
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil{	
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `uri:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccount(ctx *gin.Context){
	var req listAccountsRequest

	if err := ctx.ShouldBindQuery(req); err != nil{
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit: req.PageID-1,
		Offset: req.PageSize,
	})

	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, accounts)
}

