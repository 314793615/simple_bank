package db

import (
	"context"
	"testing"
	"github.com/314793615/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T){
	account := CreateRandomAccount(t)
	CreatRandomEntry(t, account)
	
}

func TestGetEntry(t *testing.T){
	account := CreateRandomAccount(t)
	targetEntry := CreatRandomEntry(t, account)

	entry, err := testQueries.GetEntry(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, entry.Amount, targetEntry.Amount)
	require.Equal(t, entry.CreatedAt, targetEntry.CreatedAt)
	require.Equal(t, entry.ID, targetEntry.ID)
}

func CreatRandomEntry(t *testing.T, account Account) Entry {
	createArgs := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), createArgs)
	require.NoError(t, err)
	require.Equal(t, entry.AccountID, account.ID)
	require.NotEmpty(t, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)
	return entry

}