package main

import "fmt"

type User struct {
	ID int
}

func main() {
	o1 := one()
	fmt.Println(o1)

	o2 := two()
	fmt.Println(o2)

	// Materials:
	// https://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/
	// go help build
	// go tool compile -help
	// https://blog.gopheracademy.com/code-generation-from-the-ast/
	// go build -gcflags "-S" src/cmd/alloc/alloc.go
	// GOSSAFUNC=one go build -gcflags "-S" src/cmd/alloc/alloc.go
	// https://dave.cheney.net/2018/01/08/gos-hidden-pragmas
	// https://medium.com/@ibrahimpasha.m.d/golang-toolchains-internals-intermediate-representation-df404705c32c
	// https://www.youtube.com/watch?v=2557w0qsDV0

	// Profiling
	// https://blog.golang.org/pprof
	// https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/
	// https://artem.krylysov.com/blog/2017/03/13/profiling-and-optimizing-go-web-applications/
	// https://software.intel.com/content/www/us/en/develop/blogs/debugging-performance-issues-in-go-programs.html
}

//go:noinline
func one() User {
	return User{
		ID: 1,
	}
}

//go:noinline
func two() *User {
	return &User{
		ID: 2,
	}
}
