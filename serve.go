package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

var commonIPv4 = regexp.MustCompile(`^(10|192|127)\.`)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please provide a listen address")
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("error listing interfaces: %v", err)
	}
	for _, address := range addrs {
		if commonIPv4.MatchString(address.String()) {
			log.Printf("have address %s", address)
		}
	}

	files := http.FileServer(http.Dir("."))
	log.Fatal(http.ListenAndServe(os.Args[1], http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requests %q", r.RemoteAddr, r.URL)
		files.ServeHTTP(w, r)
	})))
}
