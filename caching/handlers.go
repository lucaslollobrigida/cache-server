package caching

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	adaptor "github.com/valyala/fasthttp/fasthttpadaptor"
)

// HandleGet return the value expecified by given key, in case it is found.
func (c *Cache) HandleGet(ctx *fasthttp.RequestCtx) {
	value := c.Get(ctx.UserValue("key").(string))
	if value == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	json, err := json.Marshal(value)
	if err != nil {
		sentry.CaptureException(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetBody(json)
}

// HandlePost write the value for the expecified key.
func (c *Cache) HandlePost(ctx *fasthttp.RequestCtx) {
	c.Insert(ctx.UserValue("key").(string), string(ctx.Request.Body()))
}

// HandleDelete remove the value expecified by given key, in case it is found.
func (c *Cache) HandleDelete(ctx *fasthttp.RequestCtx) {
	c.Remove(ctx.UserValue("key").(string))
}

// Init bootstrap a router with all handlers setup.
func (c *Cache) Init() *fasthttprouter.Router {
	c.Map = make(map[string]*Registry)

	prom := adaptor.NewFastHTTPHandler(promhttp.Handler())

	router := fasthttprouter.New()

	router.GET("/health", c.HealthCheck)
	router.GET("/metrics", prom)

	router.GET("/cache/:key", c.HandleGet)
	router.DELETE("/cache/:key", c.HandleDelete)
	router.POST("/cache/:key", c.HandlePost)

	return router
}

// HealthCheck reports the health status of the server via a 200 status code.
func (c *Cache) HealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// CleanupExpired trigger a routine for cleaning expired registries from map.
func (c *Cache) CleanupExpired() {
	go func() {
		for {
			for k, v := range c.Map {
				if v.RegTime.Add(time.Minute * 5).Before(time.Now()) {
					fmt.Printf("Removed a registry: %v\n", v)
					c.Remove(k)
				}
			}
			fmt.Println("Performed a cleanup")
			time.Sleep(time.Second * 10)
		}
	}()
}
