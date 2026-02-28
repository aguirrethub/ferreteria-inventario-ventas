package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
