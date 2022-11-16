package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一次before")
				next(ctx)
				fmt.Println("第一次after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二次before")
				next(ctx)
				fmt.Println("第二次after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三次中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四次，你看不到这句话")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
