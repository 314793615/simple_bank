package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/simplebank?sslmode=disable"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(t *testing.M) {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	testQueries = New(testDB)
	os.Exit(t.Run())
}
