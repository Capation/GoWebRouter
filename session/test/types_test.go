package test

import (
	web "Go_Web"
	"Go_Web/session"
	"Go_Web/session/cookie"
	"Go_Web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {

	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(time.Minute * 15),
		CtxSessKey: "sessKey",
	}

	server := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				// 放过去，用户准备登录
				next(ctx)
				return
			}
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登录")
				return
			}

			// 刷新 session 的过期时间
			_ = m.RefreshSession(ctx)
			next(ctx)
		}
	}))

	// 登录
	server.Post("/login", func(ctx *web.Context) {
		// 要在这之前校验用户名和密码
		sess, err := m.InitSession(ctx)
		if err != nil {
			// 服务器异常
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "tzh")
		if err != nil {
			// 服务器异常
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")

		return
	})

	// 登出
	server.Post("/logout", func(ctx *web.Context) {
		// 要在这之前校验用户名和密码
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出登录")

	})

	server.Get("/user", func(ctx *web.Context) {
		sess, _ := m.GetSession(ctx)
		// 假如说我要把昵称从session中拿出来
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})

	server.Start(":8081")
}
