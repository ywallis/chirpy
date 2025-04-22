package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("HI!!!\n"))
}
func middlewareAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-auth") != "secret" {
			fmt.Println("Unauthorized attempt")
			w.WriteHeader(401)
			w.Write([]byte("Get out!\n"))
		} else {
			fmt.Println("Authorized, welcome")
			next.ServeHTTP(w, r)
		}
	})
}
func middlewareLogger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Logging request\n")
		next.ServeHTTP(w, r)
	})
}
