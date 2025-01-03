package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	mockdb "github.com/314793615/simplebank/db/mock"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestCreateAccount(t *testing.T) {
	account := RandomAccount()
	cases := []struct{
		name string
		body gin.H
		buildStud  func(t *testing.T, store *mockdb.MockStore)
		setUpAuth func(t *testing.T, server *Server, req *http.Request)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)

	}{
		{
			name: "OK",
			body: gin.H {
				"owner": account.Owner,
				"balance": account.Balance,
				"currency": account.Currency, 
			},
			buildStud: func(t *testing.T, store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner: account.Owner,
					Balance: account.Balance,
					Currency: account.Currency,
				}
				store.EXPECT().
				CreateAccount(gomock.Any(), arg).
				Return(account, nil).
				Times(1)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H {
				"owner": account.Owner,
				"balance": account.Balance,
				"currency": account.Currency, 
			},
			buildStud: func (t *testing.T, store *mockdb.MockStore)  {
				arg := db.CreateAccountParams{
					Owner: account.Owner,
					Balance: account.Balance,
					Currency: account.Currency,
				}
				store.EXPECT().
				CreateAccount(gomock.Any(), arg).
				Return(db.Account{}, sql.ErrConnDone).
				Times(1)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func (t *testing.T, resp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "badrequest",
			body: gin.H {
				"owner": account.Owner,
				"balance": account.Balance,
				"curr": account.Currency, 
			},
			buildStud: func (t *testing.T, store *mockdb.MockStore)  {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func (t *testing.T, resp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},

	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// new controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// buildstuds
			store := mockdb.NewMockStore(ctrl)
			c.buildStud(t, store)
			
			server := NewTestServer(store)
			server.SetUpRouter()
			
			data, err := json.Marshal(c.body)
			require.NoError(t, err)


			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			
			c.setUpAuth(t, server, request)

			response := httptest.NewRecorder()
			server.router.ServeHTTP(response, request)
			c.checkResponse(t, response)
			// require.Equal(t, http.StatusOK, response.Code)
		})
		
	}

}

func TestGetAccount(t *testing.T) {
	account := RandomAccount()
	cases := []struct{
		name string
		buildStud  func(t *testing.T, store *mockdb.MockStore)
		setUpAuth func(t *testing.T, server *Server, req *http.Request)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)

	}{
		{
			name: "OK",
			buildStud: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), account.ID).
				Return(account,nil).
				Times(1)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			name: "NotFound",
			buildStud: func (t *testing.T, store *mockdb.MockStore)  {
				store.EXPECT().
				GetAccount(gomock.Any(), account.ID).
				Return(db.Account{}, sql.ErrNoRows).
				Times(1)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func (t *testing.T, resp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "InternalError",
			buildStud: func (t *testing.T, store *mockdb.MockStore)  {
				store.EXPECT().
				GetAccount(gomock.Any(), account.ID).
				Return(db.Account{}, sql.ErrConnDone).
				Times(1)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				createAndAddAuth(t, server, req, account.Owner)
			},
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, resp.Code)
			},
		},
		{
			name: "NoAuth",
			buildStud: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(0)
			},
			setUpAuth:func(t *testing.T, server *Server, req *http.Request){
				// addAuth(t, server, req, account.Owner)
			},
			checkResponse: func (t *testing.T, resp *httptest.ResponseRecorder)  {
				require.Equal(t, http.StatusUnauthorized, resp.Code)
			},
		},	

	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// new controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// buildstuds
			store := mockdb.NewMockStore(ctrl)
			c.buildStud(t, store)
			
			server := NewTestServer(store)
			server.SetUpRouter()
			
			url := fmt.Sprintf("/accounts/%d", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			
			c.setUpAuth(t, server, request)

			response := httptest.NewRecorder()
			server.router.ServeHTTP(response, request)
			c.checkResponse(t, response)
			// require.Equal(t, http.StatusOK, response.Code)
		})
		
	}
	

}

func RandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 10000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}


// func checkRespBody(t *testing.T, resp *httptest.ResponseRecorder, account db.Account) {
// 	data, err := io.ReadAll(resp.Body)
// 	require.NoError(t, err)
	
// 	var gotAccount db.CreateAccountParams

// }