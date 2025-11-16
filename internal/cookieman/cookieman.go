package cookieman

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type CookieMan struct {
	key []byte
}

func newCookieMan(key []byte) *CookieMan {
	return &CookieMan{
		key,
	}
}

func NewMiddlewareCookieMan(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// handler function
		fn := func(w http.ResponseWriter, r *http.Request) {
			// new cookie manager
			cookMan := newCookieMan(key)
			// read all cookies
			cookies := r.Cookies()
			// make/validate client cookie
			clientCookie, err := cookMan.ValidateClientCookie(cookies)
			// redefine reader
			r2 := cookMan.newRequestWithContextValue(r, clientCookie, err)
			r = r2
			// set cookie to writer
			http.SetCookie(w, clientCookie)
			// response wrapper
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// next middleware handler
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (cm *CookieMan) newCookieWithErr(err error) (*http.Cookie, error) {
	cookieMsg := cm.cookieBytesToMsg(cm.makeCookie(cm.randID(), cm.key))
	return &http.Cookie{
		Name:  common.CookieUserKeyName,
		Value: cookieMsg}, err
}

func (cm *CookieMan) newRequestWithContextValue(r *http.Request, cookie *http.Cookie, err error) *http.Request {
	ctxValues := common.NewCtxValues()
	ctxValues.Set(common.CookieUserKeyName, cookie.Value)
	ctxValues.Set(common.CookieUserKeyErrName, err)
	return r.WithContext(context.WithValue(r.Context(), common.CtxKeyName, ctxValues))
}

func (cm *CookieMan) ValidateClientCookie(cookies []*http.Cookie) (*http.Cookie, error) {
	// check length of cookies
	if len(cookies) == 0 {
		return cm.newCookieWithErr(common.ErrNoCookie)
	}

	// finding needed cookie if exists
	for _, cookie := range cookies {
		// checking
		if cookie.Name == common.CookieUserKeyName {
			// convert cookie value to bytes
			b, err := cm.cookieMsgToBytes(cookie.Value)
			if err != nil {
				// error converting, generating new cookie
				return cm.newCookieWithErr(common.ErrNotValidCookie)
			}
			if _, err = cm.checkCookie(b, cm.key); err != nil {
				// cookie validation error, generating new cookie
				return cm.newCookieWithErr(common.ErrNotValidCookie)
			} else {
				return cookie, common.ErrCookieIsOk
			}
		}
	}

	// no client cookie find on slice of cookies
	return cm.newCookieWithErr(common.ErrNoCookie)
}

func (cm *CookieMan) randID() int64 {
	return common.RandID()
}

func (cm *CookieMan) checkCookie(cookie []byte, key []byte) (int64, error) {
	// get user identifier and checksum
	userID := binary.BigEndian.Uint64(cookie[:8])

	// signing with HMAC using SHA-256
	h := hmac.New(sha256.New, key)
	h.Write(cookie[:8])
	sign := h.Sum(nil)

	// check sign
	if !hmac.Equal(sign, cookie[8:]) {
		return -1, common.ErrNotValidCookie
	}

	return int64(userID), nil
}

func (cm *CookieMan) makeCookie(userID int64, key []byte) []byte {
	// convert to []byte
	byteUserID := make([]byte, 8)
	binary.BigEndian.PutUint64(byteUserID, uint64(userID))

	// signing with HMAC using SHA-256
	h := hmac.New(sha256.New, key)
	h.Write(byteUserID)
	hSum := h.Sum(nil)

	// make cookie
	cookieBytes := append(byteUserID, hSum...)

	return cookieBytes
}

func (cm *CookieMan) cookieBytesToMsg(cookieBytes []byte) string {
	return fmt.Sprintf("%x", cookieBytes)
}

func (cm *CookieMan) cookieMsgToBytes(cookieMsg string) ([]byte, error) {
	return hex.DecodeString(cookieMsg)
}
