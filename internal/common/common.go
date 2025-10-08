package common

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	RandomCharset = "abcdefghijklmnopqrstuvwxyz" + // lowercase
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + // uppercase
		"0123456789" // digits
	MinRndStrLen  = 16
	AppSrvVersion = "AVitaminOz-z-z HTTP-Server v0.3"
	CtxWaitMinSec = 10
	CtxWaitMaxSec = 30

	// MigrationScriptSep

	MigrationPgUpFile   = "migrations/000001_ya_shortener_pgstorage.up.sql"
	MigrationPgDownFile = "migrations/000001_ya_shortener_pgstorage.down.sql"

	MigrationScriptSep = "--$$--\n"
)

const (
	DefAvailableContentTypeRgx = `^text/plain(|.+)$`
	APIAvailableContentTypeRgx = `^application/json(|.+)$`
	DefContentType             = "text/plain; charset=utf-8"
	APIContentType             = "application/json"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var ErrStatusConflict = errors.New(http.StatusText(http.StatusConflict))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetRandomString(length int) string {
	return StringWithCharset(length, RandomCharset)
}

func GetSHA256Sum(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
