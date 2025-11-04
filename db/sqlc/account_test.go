package sqlc

import (
	"context"
	"testing"

	"github.com/joserafaelSH/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func tearUpAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func tearDownAccount(t *testing.T, account int64) {
	err := testStore.DeleteAccount(context.Background(), account)
	require.NoError(t, err)
}

func TestCreateAccount(t *testing.T) {
	acc := tearUpAccount(t)
	tearDownAccount(t, acc.ID)
}

func TestGetAccount(t *testing.T) {
	acc := tearUpAccount(t)
	account, err := testStore.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, acc.ID, account.ID)
	require.Equal(t, acc.Owner, account.Owner)
	require.Equal(t, acc.Balance, account.Balance)
	require.Equal(t, acc.Currency, account.Currency)
	require.WithinDuration(t, acc.CreatedAt, account.CreatedAt, 0)
	tearDownAccount(t, acc.ID)
}

func TestDeleteAccount(t *testing.T) {
	acc := tearUpAccount(t)
	tearDownAccount(t, acc.ID)
	account, err := testStore.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for range 10 {
		lastAccount = tearUpAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
	for i := 0; i < 5; i++ {
		tearDownAccount(t, accounts[i].ID)
	}
	for i := 0; i < 5; i++ {
		tearDownAccount(t, lastAccount.ID-int64(i))
	}
}

func TestUpdateAccount(t *testing.T) {
	acc := tearUpAccount(t)
	arg := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomMoney(),
	}
	account, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	//account, err := testStore.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, acc.ID, account.ID)
	require.Equal(t, acc.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, acc.Currency, account.Currency)
	require.WithinDuration(t, acc.CreatedAt, account.CreatedAt, 0)
	tearDownAccount(t, acc.ID)
}
