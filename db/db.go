package db

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var dbAddr = flag.String("db_address", "root:@(127.0.0.1:3306)/snell?", "Address of the mysql DB")
var testDbAddr = flag.String("test_db_address", "root:@(127.0.0.1:3306)/test?", "Address of test mysql DB")
var schemaPath = "/schema.sql"

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", *dbAddr)
}

func ConnectForTesting(t *testing.T) *sql.DB {
	dbc, err := sql.Open("mysql", *testDbAddr)
	if err != nil {
		t.Fatalf("db connect error: %v", err)
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	schema, err := ioutil.ReadFile(basepath + schemaPath)
	if err != nil {
		t.Errorf("Error reading schema: %s", err.Error())
		return nil
	}

	for _, q := range strings.Split(string(schema), ";") {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		q = strings.Replace(q, "create table", "create temporary table", 1)

		_, err = dbc.Exec(q)
		if err != nil {
			t.Fatalf("Error executing %s: %s", q, err.Error())
			return nil
		}
	}

	return dbc
}
