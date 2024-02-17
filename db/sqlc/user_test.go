package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yuuban007/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	username := util.RandomOwner()
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       username,
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	// how to test field PasswordChangedAt?
	// Make sure time isZero
	// require.Zero(t, user.PasswordChangedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Username, user2.Username)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

/*
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, arg.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err2 := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err2)
	require.Empty(t, account2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
}
func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testStore.ListAccount(context.Background(), arg)
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
*/
