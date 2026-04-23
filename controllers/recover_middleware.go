package controllers

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				w.WriteHeader(http.StatusInternalServerError)

				log.Printf(
					"panic recovered, method: %s, path: %s, error: %s, stack: %s\n",
					r.Method, r.URL.Path, err, string(debug.Stack()),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
