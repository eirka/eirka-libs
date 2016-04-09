package db

import (
	"database/sql"
	"fmt"
	// mysql support
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db *sql.DB

// Database holds the connection options
type Database struct {
	// Database connection settings
	User           string
	Password       string
	Proto          string
	Host           string
	Database       string
	MaxIdle        int
	MaxConnections int
}

// NewDb initializes a connection to MySQL and tries to connect.
func (d *Database) NewDb() {
	var err error

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true",
		d.User,
		d.Password,
		d.Proto,
		d.Host,
		d.Database,
	))
	if err != nil {
		panic(err)
	}

	// set max open connections
	db.SetMaxOpenConns(d.MaxConnections)
	// set max idle connections
	db.SetMaxIdleConns(d.MaxIdle)

	// try connecting to the database
	err = db.Ping()
	if err != nil {
		panic(err)
	}

}

// NewTestDb gets a database mock for testing
func NewTestDb() (mock sqlmock.Sqlmock, err error) {
	db, mock, err = sqlmock.New()
	return
}

// CloseDb closes the connection to MySQL
func CloseDb() (err error) {
	return db.Close()
}

// GetDb returns a connection to MySQL
func GetDb() (*sql.DB, error) {
	return db, nil
}

// GetTransaction will return a transaction
func GetTransaction() (*sql.Tx, error) {
	return db.Begin()
}
