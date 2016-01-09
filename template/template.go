package template

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  template.HTML
}

func MakePage(title string, body string) *Page {
	return &Page{Title: title, Body: template.HTML(body)}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("template/" + tmpl + ".html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, p)
}
