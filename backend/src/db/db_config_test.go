package db_test

import (
	"context"
	"testing"

	"github.com/ivolejon/pivo/db"
	"github.com/stretchr/testify/require"
)

func TestDbGetPool(t *testing.T) {
	_, err := db.ConnectAndGetPool(context.Background())
	require.NoError(t, err)
}

func TestDbCloseConnection(t *testing.T) {
	dbCtx, err := db.ConnectAndGetPool(context.Background())
	require.NoError(t, err)
	dbCtx.Close()
}

func TestDbPing(t *testing.T) {
	ctx := context.Background()
	dbCtx, err := db.ConnectAndGetPool(context.Background())
	require.NoError(t, err)
	err = dbCtx.Ping(ctx)
	require.NoError(t, err)
	// dbCtx.Close()
}
