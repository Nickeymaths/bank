package db

import (
	"context"
	"testing"
	"time"

	"github.com/Nickeymaths/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account *Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQuery.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	createRandomEntry(t, &account)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry1 := createRandomEntry(t, &account)
	entry2, err := testQuery.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, &account)
	}

	arg := ListEntryParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQuery.ListEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, 5, len(entries))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
