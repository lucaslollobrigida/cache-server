package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Build Ok")
	fasthttp.ListenAndServe(":3001", nil)
}
