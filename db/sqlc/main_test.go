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

var testQueries *Queries

func TestMain(t *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	testQueries = New(conn)
	os.Exit(t.Run())
}
