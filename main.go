package main

import (
	"fmt"
	"net/http"
)

func main(){

	const filePathRoot string = "."
	const port string = "8080"
	server := http.NewServeMux()

	s := http.Server{
		Addr:                         ":" + port,
		Handler:                      server,
	}

	server.Handle("/", http.FileServer(http.Dir(filePathRoot)))
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Error while starting server: %s\n", err)
	}
}
