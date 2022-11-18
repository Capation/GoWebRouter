package recover

import (
	web "Go_Web"
)

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	Log        func(ctx *web.Context)
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					// 篡改掉它
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.Data
					m.Log(ctx)
				}
			}()
			next(ctx)
		}
	}
}
