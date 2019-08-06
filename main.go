package main

import (
	"flag"
	"fmt"
	"time"

	"cache-server/caching"

	"net/http"
	_ "net/http/pprof"

	"github.com/getsentry/sentry-go"
	"github.com/valyala/fasthttp"
)

var (
	addr      = flag.String("addr", ":80", "TCP address")
	dsnSentry = flag.String("dsn", "", "Sentry DSN")
)

func main() {
	flag.Parse()

	config := sentry.ClientOptions{
		Dsn: *dsnSentry,
	}

	if err := sentry.Init(config); err != nil {
		panic(err)
	}

	defer sentry.Flush(time.Second * 5)

	c := caching.Cache{}

	router := c.Init()

	c.CleanupExpired()

	go func() {
		sentry.CaptureException(http.ListenAndServe(":6060", nil))
	}()

	fmt.Printf("Service is listen to 0.0.0.0%s\n", *addr)
	if err := fasthttp.ListenAndServe(fmt.Sprintf("0.0.0.0%s", *addr), router.Handler); err != nil {
		sentry.CaptureException(err)
	}
}
