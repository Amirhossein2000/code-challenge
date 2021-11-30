package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
)

var (
	db *sql.DB

	mysqlUsername string
	mysqlPassword string
	mysqlAddress  string
	mysqlDbName   string
)

func init() {
	mysqlUsername = os.Getenv("MYSQL_USERNAME")
	mysqlPassword = os.Getenv("MYSQL_PASSWORD")
	mysqlAddress = os.Getenv("MYSQL_ADDRESS")
	mysqlDbName = os.Getenv("MYSQL_DB_NAME")
}

func getDB() *sql.DB {
	var doOnce sync.Once
	doOnce.Do(
		func() {
			var err error
			db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
				mysqlUsername, mysqlPassword, mysqlAddress, mysqlDbName))
			if err != nil {
				logger.Errorf("connect to mysql err: %s", err.Error())
			}
		})
	return db
}
