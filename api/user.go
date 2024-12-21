package api

import (
	"net/http"
	"os/user"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct{
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type userResponse struct{
	UserName string `json:"username"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	PasswordChangeAt time.Time `json:"password_changed_at"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse{
	return userResponse{
		UserName: user.,
		FullName: user.full_name,
	}
}

func (s *Server) CreateUser(ctx *gin.Context){
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	rsp := newUserResponse()

	ctx.JSON(http.StatusOK, rsp)
}