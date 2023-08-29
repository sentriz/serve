package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please provide a listen address")
	}

	files := http.FileServer(http.Dir("."))
	log.Fatal(http.ListenAndServe(os.Args[1], http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requests %q", r.RemoteAddr, r.URL)
		files.ServeHTTP(w, r)
	})))
}
