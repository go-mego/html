package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-mego/html"
	"github.com/go-mego/mego"
)

func main() {
	e := mego.Default()
	e.Use(html.New(&html.Options{
		Directory: "./templates",
		Templates: []*html.Template{
			{
				Name: "main",
				Functions: template.FuncMap{
					"h1": func(s string) template.HTML {
						return template.HTML(fmt.Sprintf("<h1>%s</h1>", s))
					},
					"h2": func(s string) template.HTML {
						return template.HTML(fmt.Sprintf("<h2>%s</h2>", s))
					},
					"h3": func(s string) template.HTML {
						return template.HTML(fmt.Sprintf("<h3>%s</h3>", s))
					},
				},
				File: "main",
			},
		},
	}))
	e.GET("/", func(c *mego.Context, r *html.Renderer) {
		r.Render(http.StatusOK, "main")
	})
	e.Run()
}
