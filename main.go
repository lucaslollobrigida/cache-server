package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cache-server/caching"

	"net/http"
	_ "net/http/pprof"

	"github.com/getsentry/sentry-go"
	"github.com/valyala/fasthttp"
)

func main() {
	addr, ok := os.LookupEnv("PORT")
	if !ok {
		addr = ":3001"
	}

	dsnSentry, _ := os.LookupEnv("SENTRY_DSN")

	config := sentry.ClientOptions{
		Dsn: dsnSentry,
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

	go func() {
		fmt.Printf("Service is listen to 0.0.0.0%s\n", addr)
		if err := fasthttp.ListenAndServe(fmt.Sprintf("0.0.0.0%s", addr), router.Handler); err != nil {
			sentry.CaptureException(err)
		}
	}()

	waitExitSignal()
}

func waitExitSignal() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		done <- true
	}()

	<-done
}
