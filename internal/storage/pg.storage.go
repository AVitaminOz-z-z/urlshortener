package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PgDB struct {
	*sql.DB
}

func NewPgDB(dbType, connStr string) (*PgDB, error) {
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PgDB{db}, nil
}
