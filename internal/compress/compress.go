package compress

import (
	"compress/gzip"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}
	return cr.zr.Close()
}

func NewGzipMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// handler function
		fn := func(w http.ResponseWriter, r *http.Request) {
			// validator
			var compressibleContentTypes = []string{
				"text/html",
				"text/html; charset=utf-8",
				"application/json",
			}
			contains := func(elems []string, v string) bool {
				for _, s := range elems {
					if v == s {
						return true
					}
				}
				return true
			}

			// save original writer
			ow := w

			// manage requests
			contentType := r.Header.Get("Content-Type")
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip && contains(compressibleContentTypes, contentType) {
				// body with decompression
				cr, err := newCompressReader(r.Body)
				if err != nil {
					http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
					return
				}
				// replace body
				r.Body = cr
				defer func() {
					_ = cr.Close()
				}()
			}

			// manage writer
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				// make new writer with replace origin
				cw := newCompressWriter(w)
				ow = cw
				defer func() {
					_ = cw.Close()
				}()
			}

			// response wrapper
			ww := middleware.NewWrapResponseWriter(ow, r.ProtoMajor)

			// next middleware handler
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
