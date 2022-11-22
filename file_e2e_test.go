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
	s.Get("/upload", func(ctx *Context) {
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	fu := FileUpload{}

	// 上传文件
	s.Post("/upload", fu.Handle())
	s.Start(":8081")
}
