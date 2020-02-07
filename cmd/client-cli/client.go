package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
//	"IBank_core"
)
func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open data base %v",err)
	}
	defer func() {
		err := db.Close()
		log.Fatalf("can't close data base %v",err)
	}()
	err = db.Ping()
	if err != nil {
		log.Fatalf("can't ping data base %v",err)
	}

}
