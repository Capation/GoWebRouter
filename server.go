package web

import "net/http"

type HandleFunc func(ctx *Context)

type HTTPServerOption func(server *HTTPServer)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router

	// 在server层面支持middleware
	mdls []Middleware
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 最后一个是这个
	root := s.serve
	// 从后往前
	// 把后一个当作前一个的next组装成链条
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}
	// 这里执行的时候，就是从前往后了
	root(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handler(ctx)
}
