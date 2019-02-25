package caching

import (
	"encoding/json"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func (c *Cache) HandleGet(ctx *fasthttp.RequestCtx) {
	value := c.Get(ctx.UserValue("key").(string))
	if value == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	json, err := json.Marshal(value)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetBody(json)
}

func (c *Cache) HandlePost(ctx *fasthttp.RequestCtx) {
	c.Insert(ctx.UserValue("key").(string), string(ctx.Request.Body()))
}

func (c *Cache) HandleDelete(ctx *fasthttp.RequestCtx) {
	c.Remove(ctx.UserValue("key").(string))
}

func (c *Cache) Init() *fasthttprouter.Router {
	c.Map = make(map[string]*Registry)

	router := fasthttprouter.New()

	router.GET("/check", c.HealthCheck)

	router.GET("/cache/:key", c.HandleGet)
	router.DELETE("/cache/:key", c.HandleDelete)
	router.POST("/cache/:key", c.HandlePost)

	return router
}

func (c *Cache) HealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("Service is up"))
}
