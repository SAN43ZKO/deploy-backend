package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Log(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"query":    r.URL.RawQuery,
			"duration": time.Since(start),
		}).Info("request handled")
	})
}
