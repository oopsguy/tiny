package tiny

import (
	"html/template"
	"log"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodDelete  = "DELETE"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodOptions = "OPTIONS"
	MethodHead    = "HEAD"
)

type Handler func(Ctx) error

type MiddlewareHandler func(next Handler) Handler

type config struct {
	resultResolvers  map[string]ResultResolver
	templateResolver TemplateResolver
}

type Tiny struct {
	config *config
	server Server
}

func Default() *Tiny {
	c := &config{
		resultResolvers: make(map[string]ResultResolver),
	}
	s := newServer(c)
	t := &Tiny{
		config: c,
		server: s,
	}
	// register built-in view result resolvers
	t.RegisterResultResolver(ViewResultJson, new(jsonResolver))
	t.RegisterResultResolver(ViewResultContent, new(contentResolver))
	t.RegisterResultResolver(ViewResultHTML, new(htmlResolver))
	return t
}

func (t *Tiny) UseServer(server Server) {
	t.server = server
}

func (t *Tiny) Run(addr string) error {
	return t.server.Run(addr)
}

func (t *Tiny) RegisterResultResolver(view string, resolver ResultResolver) {
	t.config.resultResolvers[view] = resolver
}

func (t *Tiny) RegisterTemplateResolver(resolver TemplateResolver) {
	t.config.templateResolver = resolver
}

func (t *Tiny) SupportTemplate(glob string, funcMap template.FuncMap) {
	tpl, err := template.ParseGlob(glob)
	if err != nil {
		log.Fatalln("Parse HTML templates failed: ", err.Error())
	}
	tpl.Funcs(funcMap)
	t.RegisterTemplateResolver(&DefaultTplResolver{
		Tpl: tpl,
	})
}

func (t *Tiny) Router() Router {
	return t.server.Router()
}

func (t *Tiny) Static(pattern, folder string) {
	t.server.Static(pattern, folder)
}
