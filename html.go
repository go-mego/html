package html

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/go-mego/mego"
)

var (
	ErrNotFound = errors.New("html: the template is not found")
)

func New(option *Options) mego.HandlerFunc {
	// if len(option.Delimiter) == 2 {
	//
	// 	option.Delimiter[0]
	// 	option.Delimiter[1]
	// }

	if option.Extension == "" {
		option.Extension = "tmpl"
	}

	option.templatesMap = make(map[string]*Template)

	for _, v := range option.Templates {
		files := v.Files
		if len(files) == 0 {
			files = []string{v.File}
		}
		fn := v.Functions
		var t *template.Template

		for k, v := range files {
			path := fmt.Sprintf("%s/%s.%s", strings.TrimRight(option.Directory, "/"), strings.TrimLeft(v, "/"), option.Extension)
			content, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			if k == 0 {
				t = template.Must(template.New(v).Funcs(textFuncMap).Funcs(option.Functions).Funcs(fn).Parse(string(content)))
			} else {
				t = template.Must(t.New(v).Funcs(option.Functions).Funcs(fn).Parse(string(content)))
			}
		}

		v.templateName = files[0]
		v.template = t
		option.templatesMap[v.Name] = v
	}

	return func(c *mego.Context) {
		r := &Renderer{
			context: c,
			options: option,
		}
		c.Map(r)
		c.Next()
	}
}

type H map[string]interface{}

type Options struct {
	Directory    string
	Extension    string
	Delimiter    []string
	Functions    template.FuncMap
	Templates    []*Template
	templatesMap map[string]*Template
}

type Renderer struct {
	context *mego.Context
	options *Options
}

type Template struct {
	Name         string
	File         string
	Files        []string
	Functions    template.FuncMap
	templateName string
	template     *template.Template
}

func (r *Renderer) Render(code int, templateName string, data ...interface{}) error {
	var renderData interface{}
	if len(data) > 0 {
		renderData = data[0]
	}
	v, ok := r.options.templatesMap[templateName]
	if !ok {
		return ErrNotFound
	}
	if v.File != "" {
		err := v.template.Execute(r.context.Writer, renderData)
		if err != nil {
			return err
		}
	} else {
		err := v.template.ExecuteTemplate(r.context.Writer, v.templateName, renderData)
		if err != nil {
			return err
		}
	}
	return nil
}
