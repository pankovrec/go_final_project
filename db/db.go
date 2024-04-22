package db

import (
	"os"
	"log"
	"database/sql"
	"path/filepath"
	
	_ "modernc.org/sqlite"
)

const (
	dbName = "scheduler.db"
)

func StartDb(){
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)
	log.Println("Database not found. Creating new." + err.Error())
	var install bool
	if err != nil {
		install = true
	}
	if install == true {
		db, err := sql.Open("sqlite", dbName)
		if err != nil {
			log.Fatal("Error open db.", err)
		}
		defer db.Close()
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128)
		);`)
		if err != nil {
			log.Fatal("Error create db", err)
		
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS index_date ON scheduler (date);")
	if err != nil {
		log.Fatal("Error create index", err)
	}
	}
}