package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourcename = "postgres://root:123456@localhost/bank?sslmode=disable"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(driverName, dataSourcename)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuery = New(testDB)

	os.Exit(m.Run())
}
