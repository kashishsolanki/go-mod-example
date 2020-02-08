package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DB var for sql DB reference
var DB *sql.DB

// Open function for open db connection
func Open() (*sql.DB, error) {
	dbDriver := "mysql"
	dbName := "dt-hrms"
	dbUser := "root"
	dbPass := "admin"

	DB, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	return DB, err
}

/** Sqlx connection methods*/

// OpenSqlx function for open db connection
func OpenSqlx() (*sqlx.DB, error) {
	dbDriver := "mysql"
	dbName := "dt-hrms"
	dbUser := "root"
	dbPass := "admin"

	DB, err := sqlx.Connect(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	return DB, err
}
