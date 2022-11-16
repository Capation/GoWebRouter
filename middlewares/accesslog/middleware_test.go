package accesslog

import (
	web "Go_Web"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Builder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello world")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	req.Host = "localhost:8080"
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
