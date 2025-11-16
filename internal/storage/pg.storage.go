package storage

import (
	"context"
	"database/sql"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"time"
)

/*
func (us *URLStorage) byteAToAPIShorURL(b []byte) (*m.APIShorURL, error) {
	apiSU := &m.APIShorURL{}
	err := json.Unmarshal(b, apiSU)
	if err != nil {
		return nil, err
	}
	return apiSU, nil
}

func (us *URLStorage) returnPgDBShortURLOld(ctx context.Context, url string, prefix string) ([]byte, error) {
	db := us.getPgDB()
	row := db.QueryRowContext(ctx, "select fn__return_short_url_v2($1, $2);", url, prefix)
	short := new([]byte)
	if err := row.Scan(short); err != nil {
		return nil, err
	}
	return *short, nil
}
*/

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

func (us *URLStorage) initDB(ctx context.Context) error {
	// loading migration script from file migration/migration_file_name.ext
	script, err := os.ReadFile(common.MigrationPgUpFile)
	if err != nil {
		return err
	}

	// getting parts
	parts := strings.Split(string(script), common.MigrationScriptSep)

	// running migration steps
	db := us.getPgDB()
	tx, err := db.Begin()

	for _, s := range parts {
		if ctx.Err() != nil {
			return err
		}

		if _, err = tx.ExecContext(ctx, s); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (us *URLStorage) SetPgDB(ctx context.Context, dsn string) error {
	// setting up DB-engine
	dbType := strings.Split(dsn, ":")[0]
	pgDB, err := NewPgDB(dbType, dsn)
	if err != nil {
		return err
	}
	us.PgDB = pgDB

	// initialize DB-engine
	ctxT, cancel := context.WithTimeout(ctx, common.CtxWaitMaxSec*time.Second)
	defer cancel()

	err = us.initDB(ctxT)
	if err != nil {
		return err
	}
	us.UseDBEngine = err == nil
	return nil
}

func (us *URLStorage) returnPgDBShortURL(ctx context.Context, url string, prefix string) (string, error) {
	db := us.getPgDB()
	userKey := common.GetContextCookieUserKey(ctx)
	row := db.QueryRowContext(ctx, "select * from fn__return_short_url_on_conflict($1, $2, $3);", url, userKey, prefix)
	type Result struct {
		data     string
		conflict bool
	}
	result := &Result{}
	if err := row.Scan(&result.data, &result.conflict); err != nil {
		return "", err
	}
	if result.conflict {
		return result.data, common.ErrStatusConflict
	}
	return result.data, nil
}

func (us *URLStorage) returnPgDBBatchShortURLs(ctx context.Context, batch []byte, prefix string) ([]byte, error) {
	db := us.getPgDB()
	userKey := common.GetContextCookieUserKey(ctx)
	row := db.QueryRowContext(ctx, "select fn__return_batch_short_urls($1, $2, $3);", batch, userKey, prefix)
	batchShorts := new([]byte)
	if err := row.Scan(batchShorts); err != nil {
		return nil, err
	}
	return *batchShorts, nil
}

func (us *URLStorage) returnPgDBFullURL(ctx context.Context, short string) (string, error) {
	db := us.getPgDB()
	userKey := common.GetContextCookieUserKey(ctx)
	row := db.QueryRowContext(ctx, "select fn__return_full_url($1, $2);", short, userKey)
	url := new(string)
	if err := row.Scan(url); err != nil {
		return "", err
	}
	return *url, nil
}

func (us *URLStorage) returnPgDBUSerURLs(ctx context.Context, prefix string) ([]byte, error) {
	db := us.getPgDB()
	userKey := common.GetContextCookieUserKey(ctx)
	row := db.QueryRowContext(ctx, "select * from fn__return_user_urls($1, $2);", userKey, prefix)
	type Result struct {
		data    []byte
		isEmpty bool
	}
	result := &Result{}
	if err := row.Scan(&result.data, &result.isEmpty); err != nil {
		return nil, err
	}
	if result.isEmpty {
		return result.data, common.ErrStatusNoContent
	}
	return result.data, nil
}
