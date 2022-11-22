package web

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	type user struct {
		Name string
	}
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse("Hello, {{.Name}}")
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, &user{Name: "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse("Hello, {{.Name}}")
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]string{"Name": "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestSliceData(t *testing.T) {
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, []string{"Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestFunCall(t *testing.T) {
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
{{printf "%.2f" 1.2345}}
Hello, {{.Hello "Tom" "Jerry"}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FunCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度: 2
1.23
Hello, Tom . Jerry`, buffer.String())
}

func TestForSlice(t *testing.T) {
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := .Slice}}
{{- .}}
{{$idx}}-{{$ele}}
{{end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FunCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `a
0-a
b
1-b

`, buffer.String())
}

func TestForLoop(t *testing.T) {
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := .}}
{{- $idx}},
{{- end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, make([]int, 100))
	require.NoError(t, err)
	assert.Equal(t, `0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,
`, buffer.String())
}

func TestIfElse(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("Hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6)}}
儿童:(0, 6]
{{ else if and (gt .Age 6) (le .Age 18) }}
少年:(6, 18]
{{ else }}
我的是成人 > 18
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Age: 20})
	require.NoError(t, err)
	assert.Equal(t, `
我的是成人 > 18
`, buffer.String())
}

type FunCall struct {
	Slice []string
}

func (f FunCall) Hello(first string, last string) string {
	return fmt.Sprintf("%s . %s", first, last)
}
