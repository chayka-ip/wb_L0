package sqlstore_test

import (
	"fmt"
	"os"
	"testing"
)

var (
	databaseUrl  string
	orders_table = "orders"
	databaseTest = "db_test"
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = fmt.Sprintf("host=127.0.0.1  user=postgres password=postgres dbname=%s sslmode=disable", databaseTest)
	}
	os.Exit(m.Run())
}
