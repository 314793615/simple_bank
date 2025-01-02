package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	mockdb "github.com/314793615/simplebank/db/mock"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M){
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}

func NewTestServer(store *mockdb.MockStore) *Server {
	config := &util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration: time.Minute,
	}
	server := NewServer(config, store)
	return server
}

func createAndAddAuth(t *testing.T, server *Server, req *http.Request, user string){
	token, payload, err := server.tokenMaker.CreateToken(user, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBear, token)
	req.Header.Add(authorizationHeadKey, authorizationHeader)
}