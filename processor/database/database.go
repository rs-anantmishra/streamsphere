package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	c "github.com/rs-anantmishra/streamsphere/utils/processor/config"
)

/*
- TODO: initialize DB from scripts
- TODO: Read script files to get queries to execute
*/

var DB *sql.DB

func ConnectDB() {
	file := c.Config("DB_NAME", true)
	//if db name not provided, take default.
	if len(file) == 0 {
		file = strings.ReplaceAll(`..#database#db#streamsphere.db`, "#", string(os.PathSeparator))
	}

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
	fmt.Println("Connected to database.")
}

func CloseDB() bool {

	//Check if db isn't closed, then close db
	if err := DB.Ping(); err == nil {
		DB.Close()
	}

	// Check if the connection is successful
	if err := DB.Ping(); err != nil {
		fmt.Println("Database closed.")
		return true
	}

	return false
}
