package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Open(url string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", url)
}

func RunMigrations(con *sqlx.DB) error {
	_, err := con.Exec(`CREATE TABLE IF NOT EXISTS Project (
  id TEXT PRIMARY KEY,
  name TEXT,
  status TEXT
);`)
	if err != nil {
		return err
	}

	_, err = con.Exec(`CREATE TABLE IF NOT EXISTS Task (
  id TEXT PRIMARY KEY,
  description TEXT,
  dateGet TEXT,
  dateStart TEXT,
  dateCompleted TEXT,
  priority INTEGER,
  progress INTEGER,
  node TEXT,
  status TEXT,
  projectId TEXT
);`)
	if err != nil {
		return err
	}

	return nil
}
