package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// tmpl is the template file to be merged with the given page.
const tmpl = "template.html"

var htmlAssetPath = filepath.Join("assets", "html")

var funcMap = template.FuncMap{
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"noescapetime": func(t time.Time) template.HTML {
		return template.HTML(t.String())
	},
}

func execTemplate(page string, w http.ResponseWriter, view interface{}) {
	pageData, err := Asset(filepath.Join(htmlAssetPath, page))
	if err != nil {
		logNotFoundError(w, err)
		return
	}
	tmplData, err := Asset(filepath.Join(htmlAssetPath, tmpl))
	if err != nil {
		logNotFoundError(w, err)
		return
	}
	t := template.New("template").Funcs(funcMap)
	t, err = t.Parse(string(tmplData))
	if err != nil {
		logServerError(w, err)
		return
	}
	t, err = t.Parse(string(pageData))
	if err != nil {
		logServerError(w, err)
		return
	}
	err = t.ExecuteTemplate(w, "template", view)
	if err != nil {
		logServerError(w, err)
		return
	}
}
