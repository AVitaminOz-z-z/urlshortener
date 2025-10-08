package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func NewLogger(w io.Writer) *slog.Logger {
	//w = os.Stdout
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func NewDiscardLogger() *slog.Logger {
	//w = os.Stdout
	return slog.New(slog.DiscardHandler)
}

func NewMiddlewareLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// handler function
		fn := func(w http.ResponseWriter, r *http.Request) {
			// logging request info
			reqInfo := log.With(
				slog.String("Method", r.Method),
				slog.String("Path", r.URL.Path),
				slog.String("RemoteAddr", r.RemoteAddr),
				slog.String("UserAgent", r.UserAgent()),
			)
			// response wrapper
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// starting request time
			t1 := time.Now()
			// defer logging
			defer func() {
				reqInfo.Info("Request info",
					slog.Int("Status", ww.Status()),
					slog.Int("Bytes", ww.BytesWritten()),
					slog.String("Duration", time.Since(t1).String()),
				)
			}()
			// next middleware handler
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
