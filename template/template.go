package template

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title   string
	Links   []Link
	Request template.HTML
}

type Link struct {
	Url  string
	Desc string
}

func MakePage(title string, links []Link, request string) *Page {
	return &Page{Title: title, Links: links, Request: template.HTML(request)}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("template/" + tmpl + ".html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, p)
}
