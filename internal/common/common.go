package common

import (
	"context"
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
	AppSrvVersion = "AVitaminOz-z-z HTTP-Server v0.4"
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

const (
	CtxKeyName           = "CTX_VALUES"
	CookieUserKeyName    = "USER_KEY"
	CookieUserKeyErrName = "USER_KEY_ERR"
)

const (
	MinRange = 1
	MaxRange = 1<<63 - 1
)

var (
	seededRand         = rand.New(rand.NewSource(time.Now().UnixNano()))
	ErrStatusConflict  = errors.New(http.StatusText(http.StatusConflict))
	ErrStatusNoContent = errors.New(http.StatusText(http.StatusNoContent))
	ErrCheckCookie     = errors.New(http.StatusText(http.StatusUnauthorized))
	ErrNoCookie        = errors.New("ErrNoCookie")
	ErrNotValidCookie  = errors.New("ErrNotValidCookie")
	ErrCookieIsOk      = errors.New("ErrCookieIsOk")
)

type CtxValues struct {
	kv map[string]any
}

func NewCtxValues() *CtxValues {
	return &CtxValues{
		kv: make(map[string]any),
	}
}

func (v *CtxValues) Get(key string) any {
	return v.kv[key]
}

func (v *CtxValues) Set(key string, val any) {
	if _, ok := v.kv[key]; !ok {
		v.kv[key] = val
	}
}

func getContextUserValue(ctxKey, userKey string, ctx context.Context) any {
	var val any
	if i := ctx.Value(ctxKey); i != nil {
		if t, ok := i.(*CtxValues); ok {
			val = t.Get(userKey)
		}
	}
	return val
}

func GetContextCookieUserKey(ctx context.Context) string {
	var userKey string
	i := getContextUserValue(CtxKeyName, CookieUserKeyName, ctx)
	if i != nil {
		if t, ok := i.(string); ok {
			userKey = t
		}
	}
	return userKey
}

func GetContextCookieUserErr(ctx context.Context) error {
	var userErr error
	i := getContextUserValue(CtxKeyName, CookieUserKeyErrName, ctx)
	if i != nil {
		if t, ok := i.(error); ok {
			userErr = t
		}
	}
	return userErr
}

func randRange(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}

func RandID() int64 {
	return randRange(MinRange, MaxRange)
}

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
