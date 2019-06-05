package tiny

import (
	"path"
)

type (
	Router interface {
		Handlers() []HandlerInfo
		Module(prefix string) Router
		Use(middleware ...MiddlewareHandler) Router
		Add(pattern, method string, handler Handler, middleware ...MiddlewareHandler) Router
		ANY(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		GET(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		POST(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		DELETE(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		PUT(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		PATCH(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		OPTIONS(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
		HEAD(pattern string, handler Handler, middleware ...MiddlewareHandler) Router
	}

	HandlerInfo struct {
		Handler    Handler
		Method     string
		Pattern    string
		Middleware []MiddlewareHandler
	}

	defaultRouter struct {
		Router
		prefix     string
		modules    []*defaultRouter
		handlers   []HandlerInfo
		middleware []MiddlewareHandler
	}
)

func (r *defaultRouter) Handlers() []HandlerInfo {
	return r.flatHandlers("")
}

func (s *defaultRouter) Use(middleware ...MiddlewareHandler) Router {
	s.middleware = append(s.middleware, middleware...)
	return s
}

func (r *defaultRouter) Module(prefix string) Router {
	m := &defaultRouter{
		prefix: prefix,
	}
	r.modules = append(r.modules, m)
	return m
}

func (r *defaultRouter) Add(pattern, method string, handler Handler, middleware ...MiddlewareHandler) Router {
	if pattern == "" {
		pattern = "/"
	}
	r.handlers = append(r.handlers, HandlerInfo{
		Pattern:    pattern,
		Method:     method,
		Handler:    handler,
		Middleware: middleware,
	})
	return r
}

func (r *defaultRouter) ANY(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, "", handler, middleware...)
}

func (r *defaultRouter) GET(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodGet, handler, middleware...)
}

func (r *defaultRouter) POST(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodPost, handler, middleware...)
}

func (r *defaultRouter) PUT(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodPut, handler, middleware...)
}

func (r *defaultRouter) DELETE(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodDelete, handler, middleware...)
}

func (r *defaultRouter) PATCH(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodPatch, handler, middleware...)
}

func (r *defaultRouter) OPTIONS(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodOptions, handler, middleware...)
}

func (r *defaultRouter) HEAD(pattern string, handler Handler, middleware ...MiddlewareHandler) Router {
	return r.Add(pattern, MethodHead, handler, middleware...)
}

func (r *defaultRouter) flatHandlers(prefix string, parentMiddleware ...MiddlewareHandler) []HandlerInfo {
	var handlers []HandlerInfo
	for _, h := range r.handlers {
		h.Pattern = path.Join(prefix, r.prefix, h.Pattern)
		h.Middleware = append(append(parentMiddleware, r.middleware...), h.Middleware...)
		handlers = append(handlers, h)
	}
	for _, m := range r.modules {
		// handlers = append(handlers, flatModule(r.prefix, m, r.middleware...)...)
		handlers = append(
			handlers,
			m.flatHandlers(
				path.Join(prefix, r.prefix),
				append(parentMiddleware, r.middleware...)...,
			)...,
		)
	}
	return handlers
}
