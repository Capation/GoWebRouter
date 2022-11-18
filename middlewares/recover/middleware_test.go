package recover

import (
	web "Go_Web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {

	builder := MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("你 panic 了"),
		Log: func(ctx *web.Context) {
			fmt.Printf("panic 的路径是: %s", ctx.Req.URL.String())
		},
	}

	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		panic("发生panic了")
	})
	server.Start(":8083")
}
