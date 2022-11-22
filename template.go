package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// tplName 模板的名字，按名索引
	// data 渲染页面用的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(buffer, tplName, data)
	return buffer.Bytes(), err
}
