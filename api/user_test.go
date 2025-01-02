package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/314793615/simplebank/db/mock"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestCreateUser(t *testing.T){
	user, password := RandomUser(t)
	cases := []struct{
		name string
		body gin.H
		buildStud func(store *mockdb.MockStore)
		// setUpAuth func(t *testing.T, server *Server, request *http.Request)
		checkResponse func(rsp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H {
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStud: func (store *mockdb.MockStore)  {
				store.EXPECT().
				CreatUser(gomock.Any(), gomock.Any()).
				Return(user, nil).
				Times(1)
			},
			// setUpAuth: func (t *testing.T, server *Server, request *http.Request)  {
			// 	createAndAddAuth(t, server, request, user.Username)
			// },
			checkResponse: func(rsp *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, rsp.Code)
				verifyResponseBody(t, user, rsp)
			},
		},
	}
	for i := range cases{
		
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			server := NewTestServer(store)
			server.SetUpRouter()
			c.buildStud(store)

			data, err := json.Marshal(c.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			recoder := httptest.NewRecorder()
			server.router.ServeHTTP(recoder, request)
			c.checkResponse(recoder)
		})
		
	}

}

func TestGetUser(t *testing.T) {
	user, _ := RandomUser(t)
	cases := []struct{
		name string
		username string
		buildStud func(store *mockdb.MockStore)
		setUpAuth func (t *testing.T, server *Server, request *http.Request)
		checkResponse func(t *testing.T, rsp *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			username: user.Username,
			buildStud: func (store *mockdb.MockStore)  {
				// arg := db.CreatUserParams{
				// 	Username:user.Username,
				// 	HashedPassword: user.HashedPassword,
				// 	FullName: user.FullName,
				// 	Email: user.Email,
				// }
				store.EXPECT().
				GetUser(gomock.Any(), user.Username).
				Return(user, nil).
				Times(1)
			},
			setUpAuth: func (t *testing.T, server *Server, request *http.Request)  {
				createAndAddAuth(t, server, request, user.Username)
			},
			checkResponse: func(t *testing.T, rsp *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, rsp.Code)
				verifyResponseBody(t, user, rsp)
			},
		},

	}
	for i := range cases{
		
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			server := NewTestServer(store)
			server.SetUpRouter()
			c.buildStud(store)

			url := fmt.Sprintf("/users/%s", c.username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			c.setUpAuth(t, server, request)

			
			recoder := httptest.NewRecorder()
			server.router.ServeHTTP(recoder, request)
			c.checkResponse(t, recoder)
		})
		
	}

}


func RandomUser(t *testing.T) (db.User, string) {
	// username := util.RandomString(6)
	password := util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	return db.User{
		Username: util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}, password
}

func verifyResponseBody(t *testing.T, user db.User, rsp *httptest.ResponseRecorder){
	var gotUser db.User
	data, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username,gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}