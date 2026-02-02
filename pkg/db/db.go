package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title TEXT,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	if dbFile == "" {
		return fmt.Errorf("DB empty way")
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		f, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("DB create error: %w", err)
		}
		f.Close()
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("DB open error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("DB init error: %w", err)
	}

	return nil
}
