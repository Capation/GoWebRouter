//go:build e2e

package prometheus

import (
	web "Go_Web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {

	builder := MiddlewareBuilder{
		Namespace: "Gopher",
		Subsystem: "web",
		Name:      "http_response",
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))

	server.Get("/user", func(ctx *web.Context) {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)
		ctx.RespJSON(202, &User{
			Name: "Tom",
		})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()

	server.Start(":8081")
}

type User struct {
	Name string
}
