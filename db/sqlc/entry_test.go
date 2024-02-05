package db

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: int64(rand.Intn(20)),
		Amount:    int64(rand.Intn(2000) - 1000), // -1000~1000
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestListEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	arg := ListEntryParams{
		AccountID: entry1.AccountID,
		Limit:     10,
		Offset:    0,
	}
	entries, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry1.AccountID, entry.AccountID)
	}
}

func TestGetEntryById(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntryById(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
}
