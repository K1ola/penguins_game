package middleware

import (
	"game/helpers"
	"net/http"

)

type Status struct {
	http.ResponseWriter
	Code int
}

var SECRET = []byte("myawesomesecret")

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//origin := r.Header.Get("Origin")

		responseHeader := w.Header()
		responseHeader.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		responseHeader.Set("Access-Control-Allow-Credentials", "true")
		responseHeader.Set("Access-Control-Allow-Headers", "Content-Type")
		responseHeader.Set("Access-Control-Allow-Origin", "https://penguin-wars.sytes.pro")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				helpers.LogMsg("Recovered panic: ", err)
				http.Error(w, "Server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

