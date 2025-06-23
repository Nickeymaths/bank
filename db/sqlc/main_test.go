package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Nickeymaths/bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load testing config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot create connection pool")
	}
	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
