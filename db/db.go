package db

import (
	"database/sql"
	"flag"
)

var dbAddress = flag.String("db_address", "root:@(127.0.0.1:3306)/snell", "Address of the mysql DB")

func Connect() (*sql.db, error) {
	return sql.Open("mysql", *dbAddress)
}
