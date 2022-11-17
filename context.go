package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req *http.Request

	// Resp 如果用户直接使用这个
	// 那么用户就绕开了 RespStatusCode 和 RespData 这两个
	// 那么部分的middleware就无法运作
	Resp http.ResponseWriter

	// 这个主要是为了 middleware 读写用的
	RespStatusCode int
	RespData       []byte

	PathParams map[string]string

	queryValues url.Values

	MatchedRoute string
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	//c.Resp.WriteHeader(status)
	//n, err := c.Resp.Write(data)
	//if n != len(data) {
	//	return errors.New("web: 未写入全部数据")
	//}
	c.RespData = data
	c.RespStatusCode = status
	return nil
}

func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body为空")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.Req.FormValue(key), nil
}

func (c *Context) QueryValue(key string) (string, error) {
	// 用户是区别不出来真的有值，恰好值是空字符串
	// 缓存住这个values
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	vals, ok := c.queryValues[key]
	if !ok {
		return "", errors.New("web:没有找到所对应的key")
	}
	return vals[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}
	return val, nil
}

type StringValue struct {
	val string
	err error
}

func (c *Context) PathValue1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			val: "",
			err: errors.New("web: key不存在"),
		}
	}
	return StringValue{
		val: val,
		err: nil,
	}
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
