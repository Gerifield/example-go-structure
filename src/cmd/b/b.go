package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gerifield/example-go-structure/src/example1"
)

func main() {
	listen := flag.String("listen", ":8080", "Http listen address")
	static := flag.String("static", "./static", "Static folder")
	flag.Parse()

	s := example1.New(*static)

	log.Println("Listen on", *listen)
	_ = http.ListenAndServe(*listen, s.Routes())
}
