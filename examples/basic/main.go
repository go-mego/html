package main

import (
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
				File: "main.tmpl",
			},
		},
	}))
	e.GET("/", func(c *mego.Context, r *html.Renderer) {
		r.Render(http.StatusOK, "main", html.H{
			"Name": "小安",
			"List": []string{"第一個", "第二個", "第三個"},
		})
	})
	e.Run()
}
