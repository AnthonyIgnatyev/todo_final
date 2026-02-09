package db

import (
	"database/sql"
	"log"
	"os"

	cfg "todo_final/pkg/config"

	_ "modernc.org/sqlite"
)

const schema string = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR NOT NULL DEFAULT ""
);
CREATE INDEX idx_date ON scheduler (date);
`

var db *sql.DB

func Init() error {
	var install bool
	dbFile := cfg.CfgStruct.Database.FilePath

	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return err
		}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
		log.Println("Database created with schema:", dbFile)
	} else {
		log.Println("Database already exists:", dbFile)
	}

	return db.Ping()
}

func CloseDB() error {
	return db.Close()
}
