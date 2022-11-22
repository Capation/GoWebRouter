//go:build e2e

package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"testing"
)

func TestLoginPage(t *testing.T) {

	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}
	s := NewHTTPServer(ServerWithTemplateEngine(engine))
	s.Get("/login", func(ctx *Context) {
		err = ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	s.Start(":8081")
}
