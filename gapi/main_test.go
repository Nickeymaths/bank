package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/token"
	"github.com/Nickeymaths/bank/util"
	"github.com/Nickeymaths/bank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func NewTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		SymmetricKey:        util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributor)

	require.NoError(t, err)
	require.NotEmpty(t, server)
	return server
}

func newContextWithBearerTokenfunc(t *testing.T, username string, role string, duration time.Duration, tokenMarker token.Maker) context.Context {
	accessToken, _, err := tokenMarker.CreateToken(username, role, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", bearerHeader, accessToken)
	mtd := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), mtd)
}
