package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	flag.Parse()

	addr := flag.Arg(0)
	if addr == "" {
		log.Fatal("need listen addr")
	}

	var path = "."
	if p := flag.Arg(1); p != "" {
		path = p
	}

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	var hander http.Handler
	switch {
	case info.Mode().IsRegular():
		hander = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path)
		})
	default:
		hander = http.FileServer(http.Dir(path))
	}

	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requests %q", r.RemoteAddr, r.URL)
		hander.ServeHTTP(w, r)
	})))
}
