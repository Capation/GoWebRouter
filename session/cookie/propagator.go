package cookie

import (
	"net/http"
)

type PropagatorOption func(propagator *Propagator)

type Propagator struct {
	cookieName   string
	cookieOption func(c *http.Cookie)
}

func NewPropagator() *Propagator {
	return &Propagator{
		cookieName: "sessid",
		cookieOption: func(c *http.Cookie) {

		},
	}
}

func WithCookieName(name string) PropagatorOption {
	return func(propagator *Propagator) {
		propagator.cookieName = name
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name: p.cookieName,
		// value 是 id
		Value:    id,
		HttpOnly: true,
	}
	p.cookieOption(c)
	http.SetCookie(writer, c)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	c, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// 重新设置一下，就会覆盖掉之前的
	c := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, c)
	return nil
}
