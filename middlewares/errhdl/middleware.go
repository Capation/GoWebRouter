package errhdl

import web "Go_Web"

type MiddlewareBuilder struct {
	// 不能动态渲染
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {

	// 初始化这个map
	return &MiddlewareBuilder{resp: map[int][]byte{}}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				// 篡改结果
				ctx.RespData = resp
			}
			next(ctx)
		}
	}
}
