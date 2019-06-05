package tiny

import (
	"fmt"
	"log"
	"net/http"
)

type Server interface {
	Run(addr string) error
	Static(pattern, folder string)
	Router() Router
	Shutdown() error
}

type defaultServer struct {
	Server
	config *config
	router Router
	server *http.Server
}

func newServer(c *config) Server {
	return &defaultServer{
		config: c,
		router: new(defaultRouter),
		server: new(http.Server),
	}
}

func (s *defaultServer) Static(pattern, folder string) {
	http.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(folder))))
}

func (s *defaultServer) handle(method string, pattern string, h Handler) {
	log.Printf("Registering handler: %s[%s]\n", pattern, method)
	http.HandleFunc(pattern, func(resp http.ResponseWriter, req *http.Request) {
		ctx := &defaultCtx{
			resp:   resp,
			req:    req,
			config: s.config,
		}
		// check method
		if method != "" && method != req.Method {
			ctx.Status(http.StatusMethodNotAllowed)
			_ = ctx.Content(fmt.Sprintf("%s %s Method not allowed", req.Method, req.RequestURI))
			return
		}
		log.Printf("Request: %s[%s] → handler[%s]\n", ctx.Req().RequestURI, ctx.Req().Method, pattern)
		if err := h(ctx); err != nil {
			log.Println(err)
			// TODO html, json or jsonp ?
			ctx.Status(http.StatusInternalServerError)
			_ = ctx.Content(err.Error())
			return
		}
	})
}

func (s *defaultServer) Router() Router {
	return s.router
}

func (s *defaultServer) Run(addr string) error {
	for _, h := range s.router.Handlers() {
		for i := len(h.Middleware) - 1; i >= 0; i-- {
			h.Handler = h.Middleware[i](h.Handler)
		}
		s.handle(h.Method, h.Pattern, h.Handler)
	}
	s.server.Addr = addr
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *defaultServer) Shutdown() error {
	if err := s.server.Close(); err != nil {
		return err
	}
	return nil
}
