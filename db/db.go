package db

import (
	"os"
	"log"
	"database/sql"
	"path/filepath"
	//"fmt"
	
	_ "modernc.org/sqlite"
)

const (
	dbName = "scheduler.db"
	SQL_CREATE_TABLES                     = "CREATE TABLE IF NOT EXISTS scheduler " +
	"(id INTEGER PRIMARY KEY AUTOINCREMENT, " +
	"date TEXT, " +
	"title TEXT, " +
	"comment TEXT, " +
	"repeat VARCHAR(128));"
)

func StartDb(){
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)
	log.Println("Ошибка. БД не существует. Делаем новую" + err.Error())
	var install bool
	if err != nil {
		install = true
	}
	if install == true {
		db, err := sql.Open("sqlite", dbName)
		if err != nil {
			log.Fatal("FAILED_TO_OPEN_DATABASE", err)
		}
		defer db.Close()
		_, err = db.Exec(SQL_CREATE_TABLES);
		if err != nil {
			log.Fatal("TABLE_CREATION_ERROR", err)
		
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS index_date ON scheduler (date);")
	if err != nil {
		log.Fatal("INDEX_CREATION_ERROR", err)
	}
	}
}