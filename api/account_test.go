package api

import (
	"fmt"
	mockdb "github.com/314793615/simplebank/db/mock"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/314793615/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestCreateAccount(t *testing.T) {
	account := RandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)

	server := NewServer(&util.Config{}, store)
	server.SetUpRouter()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	server.router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)

}

func RandomAccount() db.Account {
	return db.Account{
		//ID:       util.RandomInt(1, 10000),
		ID:       3,
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
