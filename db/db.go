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

// Check a bool in the database
func GetBool(column, table, row string, id uint) (boolean bool) {

	// Check if thread is closed and get the total amount of posts
	err := db.QueryRow("SELECT ? FROM ? WHERE ? = ?", column, table, row, id).Scan(&boolean)
	if err != nil {
		return false
	}

	return

}

// Set a bool in the database
func SetBool(table, column, row string, boolean bool, id uint) (err error) {

	ps, err := db.Prepare("UPDATE ? SET ?=? WHERE ?=?")
	if err != nil {
		return
	}
	defer ps.Close()

	_, err = ps.Exec(table, column, boolean, row, id)
	if err != nil {
		return
	}

	return

}
