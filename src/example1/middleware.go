package example1

import (
	"log"
	"net/http"
	"time"
)

func authMiddleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			return // Stop it
		}

		if u == "admin" && p == "admin" {
			h.ServeHTTP(w, r)
			return
		}
		return
	})
}

type userPass struct {
	user string
	pass string
}

func authMiddleware2(users []userPass) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				return // Stop it
			}

			for _, up := range users {
				if up.user == u && up.pass == p {
					h.ServeHTTP(w, r)
					return
				}
			}

			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Nope!"))
			return
		})
	}
}

func basicLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s (%s) \n", r.Method, r.RequestURI, time.Now().Sub(t))
	})
}
