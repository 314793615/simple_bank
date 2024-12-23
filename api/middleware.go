package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/314793615/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeadKey = "authorization"
	authorizationTypeBear = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(maker token.TokenMaker) gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		authorizationHeader := ctx.GetHeader(authorizationHeadKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("authorization header is not correct")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBear {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err !=  nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return 
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}

}