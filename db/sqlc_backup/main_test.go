package db

import (
	"context"
	"log"
	"os"
	util "simplebank/utils"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://postgres:admin@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
    config, fail := util.LoadConfig(".")
	if fail != nil {
		log.Fatal("cannot load config:", fail)
	}

    var err error
    testDB, err = pgxpool.New(context.Background(), config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    testQueries = New(testDB)

    os.Exit(m.Run())
}
