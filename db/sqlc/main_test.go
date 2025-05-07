package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Nickeymaths/bank/util"
	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load testing config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQuery = New(testDB)

	os.Exit(m.Run())
}
