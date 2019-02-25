package main

import (
	"flag"
	"fmt"

	"github.com/lucaslollobrigida/cache-server/cache"
	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":3001", "TCP address")
)

func main() {
	flag.Parse()

	c := cache.Cache{}

	router := c.Init()

	fmt.Printf("Service is listen to 0.0.0.0%s\n", *addr)
	if err := fasthttp.ListenAndServe(*addr, router.Handler); err != nil {
		fmt.Println(err)
	}
}
