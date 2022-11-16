//go:build e2e

package accesslog

import (
	web "Go_Web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Builder_e2e(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello world"))
	})
	//req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	//req.Host = "localhost:8080"
	//if err != nil {
	//	t.Fatal(err)
	//}
	//server.ServeHTTP(nil, req)
	server.Start(":8081")
}
