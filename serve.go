package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please provide a listen address")
	}
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		log.Printf("have address %s", address)
	}
	files := http.FileServer(http.Dir("."))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requests %q", r.RemoteAddr, r.URL)
		files.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(os.Args[1], handler))
}
