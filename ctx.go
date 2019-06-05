package tiny

import (
	"errors"
	"net/http"
)

var (
	errViewResolverNotFound = errors.New("can not found any view resolver")
)

type Ctx interface {
	Resp() http.ResponseWriter
	Req() *http.Request
	Form(key string) string
	Forms(key string) []string
	Query(key string) string
	Queries(key string) []string
	JSON(data interface{}) error
	JSONP(data interface{}, callback string) error
	Content(content string) error
	HTML(html string) error
	Template(name string, data interface{}) error
	Status(code int)
}

type defaultCtx struct {
	Ctx
	config *config
	resp   http.ResponseWriter
	req    *http.Request
}

func (c *defaultCtx) Resp() http.ResponseWriter {
	return c.resp
}

func (c *defaultCtx) Req() *http.Request {
	return c.req
}

func (c *defaultCtx) Form(key string) string {
	s := c.Req().Form[key]
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func (c *defaultCtx) Forms(key string) []string {
	return c.Req().Form[key]
}

func (c *defaultCtx) Query(key string) string {
	q := c.Req().URL.Query()[key]
	if len(q) > 0 {
		return q[0]
	}
	return ""
}

func (c *defaultCtx) Queries(key string) []string {
	return c.Req().URL.Query()[key]
}

func (c *defaultCtx) JSON(data interface{}) error {
	return c.writeResult(data, ViewResultJson)
}

func (c *defaultCtx) JSONP(data interface{}, callback string) error {
	return c.writeResult(data, ViewResultJsonp)
}

func (c *defaultCtx) Content(content string) error {
	return c.writeResult(content, ViewResultContent)
}

func (c *defaultCtx) HTML(html string) error {
	return c.writeResult(html, ViewResultHTML)
}

func (c *defaultCtx) Template(name string, data interface{}) error {
	if c.config.templateResolver == nil {
		return errViewResolverNotFound
	}
	var err error
	var bs []byte
	bs, err = c.config.templateResolver.Render(name, data)
	if err != nil {
		return err
	}
	_, err = c.Resp().Write(bs)
	return err
}

func (c *defaultCtx) Status(code int) {
	c.Resp().WriteHeader(code)
}

func (c *defaultCtx) writeResult(data interface{}, viewType string) error {
	viewResolver := c.config.resultResolvers[viewType]
	if viewResolver == nil {
		return errViewResolverNotFound
	}
	switch viewType {
	case ViewResultJson, ViewResultJsonp:
		c.Resp().Header().Add("content-Type", "application/json;charset=UTF-8")
	case ViewResultContent:
		c.Resp().Header().Add("content-type", "text/plain;charset=UTF-8")
	case ViewResultHTML:
		c.Resp().Header().Add("content-type", "text/html;charset=UTF-8")
	default:
		return errViewResolverNotFound
	}
	var err error
	var b []byte
	if b, err = viewResolver.Render(data); err != nil {
		return err
	}
	if b == nil || len(b) == 0 {
		return nil
	}
	_, err = c.Resp().Write(b)
	return err
}
