package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gerifield/example-go-structure/src/example2"
)

func main() {
	listen := flag.String("listen", ":8080", "Http listen address")
	flag.Parse()

	s := example2.NewHTTP(example2.NewApplication())

	// pprof enable
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Listen on", *listen)
	_ = http.ListenAndServe(*listen, s.Routes())
}
