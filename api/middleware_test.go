package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/314793615/simplebank/token"
	"github.com/314793615/simplebank/util"

	mockdb "github.com/314793615/simplebank/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestMiddleWare(t *testing.T){
	cases := []struct{
		name string
		setUpAuth func(t *testing.T, request *http.Request, server *Server)
		buildStud func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rsp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func (t *testing.T, request *http.Request, server *Server)  {
				addAuthorization(t, server.tokenMaker, "bear", authorizationHeadKey, util.RandomOwner(), time.Minute, request)
			},
			checkResponse: func (t *testing.T, rsp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusOK, rsp.Code)
			},
		},
		{
			name: "NoAuth",
			setUpAuth: func (t *testing.T, request *http.Request, server *Server)  {
				
			},
			checkResponse: func (t *testing.T, rsp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusUnauthorized, rsp.Code)
			},
		},
		{
			name: "badAuthType",
			setUpAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, server.tokenMaker, "bad", authorizationHeadKey, util.RandomOwner(), time.Minute, request)
			},
			checkResponse: func(t *testing.T, rsp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rsp.Code)
			},
		},
		{
			name: "badAuthKey",
			setUpAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, server.tokenMaker, "bear", "auth", util.RandomOwner(), time.Minute, request)
			},
			checkResponse: func(t *testing.T, rsp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rsp.Code)
			},
		},
		{
			name: "invalidAuthKeyType",
			setUpAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, server.tokenMaker, "", authorizationHeadKey, util.RandomOwner(), time.Minute, request)
			},
			checkResponse: func(t *testing.T, rsp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rsp.Code)
			},
		},
		{
			name: "ExpiredToken",
			setUpAuth: func(t *testing.T, request *http.Request, server *Server) {
				addAuthorization(t, server.tokenMaker, "bear", authorizationHeadKey, util.RandomOwner(), -time.Minute, request)
			},
			checkResponse: func(t *testing.T, rsp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rsp.Code)
			},
		},

	} 

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func (t *testing.T)  {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			server := NewTestServer(store)
			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})
			recoder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recoder, request)
		})

	}
}

func addAuthorization(
	t *testing.T,
	mk token.TokenMaker,
	keyType string,
	authorizationHeaderKey string,
	username string,
	duration time.Duration,
	request *http.Request,

){
	tk, _, err := mk.CreateToken(username, duration)
	require.NoError(t, err)
	authString := fmt.Sprintf("%s %s", keyType, tk)
	request.Header.Add(authorizationHeaderKey, authString)
}



