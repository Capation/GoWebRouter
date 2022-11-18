package errhdl

import (
	web "Go_Web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {

	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
    <body>
        <h1>Not Found</h1>
    </body>
</html>
`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
    <body>
        <h1>Bad Request</h1>
    </body>
</html>
`))

	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Start(":8081")
}
