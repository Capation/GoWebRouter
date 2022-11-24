//go:build e2e

package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestPostPage(t *testing.T) {

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

	fu := FileUpload{
		// <input type="file" name="myfile" />
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", header.Filename)
		},
	}

	// 上传文件
	s.Post("/upload", fu.Handle())
	s.Start(":8081")
}

func TestDownload(t *testing.T) {

	s := NewHTTPServer()

	fu := FileDownloader{
		Dir: filepath.Join("testdata", "download"),
	}

	s.Get("/download", fu.Handle())
	s.Start(":8081")
}
