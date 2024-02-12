package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/yuuban007/simplebank/util"
)

var testStore Store
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	var config util.Config

	config, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configuration", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
