package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(d.MaxIdle)
	db.SetMaxOpenConns(d.MaxConnections)
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
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	return tx, err
}
