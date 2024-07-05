package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	fmt.Println("Connecting to database...")
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	// username:password@tcp(127.0.0.1:3306)/test
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	db = conn
	fmt.Println("connected to database...")
}

func GetDBClient() *sql.DB {
	return db
}
