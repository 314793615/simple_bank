package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type renewAccessTokenParams struct {
	refreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct{
	accessToken string `json:"access_token"`
	accessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenParams
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}
	
	refreshPayload, err := server.tokenMaker.VerifyToken(req.refreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid refresh token: %w", err)))
		return 
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	if session.Username != refreshPayload.Usename {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("the user of token is not consisten with the user in sesssion")))
		return 
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("the refresh token is expired")))
	}

	token, accessPayLoad,  err := server.tokenMaker.CreateToken(refreshPayload.Usename, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	rsp := renewAccessTokenResponse{
		accessToken: token,
		accessTokenExpiredAt: accessPayLoad.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)

}