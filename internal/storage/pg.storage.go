package storage

import (
	"context"
	"database/sql"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	_ "github.com/lib/pq"
	"strings"
	"time"
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

func (pdb *PgDB) Ping() error {
	return pdb.DB.Ping()
}

func (us *URLStorage) getPgDB() *PgDB {
	return us.PgDB
}

func (us *URLStorage) initDB() error {
	/*data, err := os.ReadFile(common.MigrationPgUpFile)
	if err != nil {
		return err
	}*/

	// loading migration script
	script := common.MigrationScriptUp

	// getting parts
	parts := strings.Split(script, common.MigrationScriptSep)

	ctx, cancel := context.WithTimeout(context.Background(), common.CtxWaitMaxSec*time.Second)
	defer cancel()

	// running migration steps
	var err error
	db := us.getPgDB()
	for _, s := range parts {
		if _, err = db.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (us *URLStorage) SetPgDB(dsn string) error {
	// setting up DB-engine
	// dbType := strings.Split(dsn, ":")[0]
	pgDB, err := NewPgDB("postgres", dsn)
	if err != nil {
		return err
	}
	us.PgDB = pgDB
	// initialize DB-engine
	err = us.initDB()
	if err != nil {
		return err
	}
	us.UseDBEngine = err == nil
	return nil
}

func (us *URLStorage) returnPgDBShortURL(url string) string {
	db := us.getPgDB()
	ctx, cancel := context.WithTimeout(context.Background(), common.CtxWaitMinSec*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "select fn__return_short_url($1);", url)
	short := new(string)
	if err := row.Scan(short); err != nil {
		panic(err)
	}
	return *short
}

func (us *URLStorage) returnPgDBFullURL(short string) string {
	db := us.getPgDB()
	ctx, cancel := context.WithTimeout(context.Background(), common.CtxWaitMinSec*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "select fn__return_full_url($1);", short)
	url := new(string)
	if err := row.Scan(url); err != nil {
		panic(err)
	}
	return *url
}
