package main

import (
	"flag"
	"net/http"

	"github.com/gerifield/example-go-structure/src/example1"
)

func main() {
	listen := flag.String("listen", ":8080", "Http listen address")
	flag.Parse()

	s := example1.New()
	_ = http.ListenAndServe(*listen, s.Routes())
}
