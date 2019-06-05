package tiny

import (
	"bytes"
	"encoding/json"
	"html/template"
)

const (
	ViewResultJson    = "json"
	ViewResultJsonp   = "jsonp"
	ViewResultContent = "content"
	ViewResultHTML    = "html"
)

type ResultResolver interface {
	Render(data interface{}) ([]byte, error)
}

type jsonResolver struct {
	ResultResolver
}

func (r *jsonResolver) Render(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

type DefaultJSONPResolver struct {
	ResultResolver
	Callback string
}

func (r *DefaultJSONPResolver) Render(data interface{}) ([]byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return []byte(r.Callback + "(" + string(d) + ")"), nil
}

type contentResolver struct {
	ResultResolver
}

func (r *contentResolver) Render(data interface{}) ([]byte, error) {
	//c := template.HTMLEscapeString(data.(string))
	return []byte(data.(string)), nil
}

type htmlResolver struct {
	ResultResolver
}

func (r *htmlResolver) Render(data interface{}) ([]byte, error) {
	return []byte(template.HTML(data.(string))), nil
}

type TemplateResolver interface {
	Render(name string, data interface{}) ([]byte, error)
}

type DefaultTplResolver struct {
	TemplateResolver
	Tpl *template.Template
}

func (r *DefaultTplResolver) Render(name string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.Tpl.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
