package start

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Isotere/logs/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Info(fmt.Sprintf("[API Query start] Method: \"%s\" URI: \"%s\" Referer: \"%s\"", r.Method, r.RequestURI, r.Referer()))

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			spent := time.Since(start).String()
			log.Info(fmt.Sprintf("[API Query finished] Method: \"%s\" URI: \"%s\" Code: \"%d\" time spent: \"%s\"", r.Method, r.RequestURI, ww.Status(), spent))
		}

		return http.HandlerFunc(fn)
	}
}
