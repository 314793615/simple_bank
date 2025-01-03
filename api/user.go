package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	UserName         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_changed_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func newUserResponse(user *db.User) userResponse {
	return userResponse{
		UserName:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := s.store.CreatUser(ctx, db.CreatUserParams{
		Username:       req.UserName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		FullName:       req.FullName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(&user)

	ctx.JSON(http.StatusOK, rsp)
}

type GetUserRequest struct {
	Username string `uri:"username"`
}

func (s *Server) GetUser(ctx *gin.Context) {
	var req GetUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(&user)

	ctx.JSON(http.StatusOK, rsp)

}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	User                  userResponse `json:"user"`
	SessionId             uuid.UUID  `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("the password is not correct")))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayLoad, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayLoad.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshPayLoad.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("cannot create session: %w", err)))
	}
	rsp := LoginUserResponse{
		User:                 newUserResponse(&user),
		SessionId: session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiredAt: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiredAt: refreshPayLoad.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)

}
