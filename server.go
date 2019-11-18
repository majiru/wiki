package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/edit/"):
		editHandler(w, r)
		return
	case r.URL.Path == "/search":
		searchHandler(w, r)
		return
	case pageFinder.MatchString(strings.TrimPrefix(r.URL.Path, "/")):
		pageHandler(w, r)
		return
	}
	http.Error(w, "Page not Found", 404)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println("searchHandler: malformed post bailing")
			http.Error(w, "Bad Request", 400)
			return
		}
		query, ok := r.Form["query"]
		if !ok || len(query) < 1 {
			log.Println("searchHandler: malformed post bailing")
			http.Error(w, "Bad Request", 400)
			return
		}
		err = MarshalSearch(w, query[0])
		if err != nil {
			http.Error(w, "Bad Request", 400)
			log.Println("searchHandler: Error in marshalling page:", err)
		}
	default:
		http.Error(w, "Bad Request", 400)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fpath := path.Base(r.URL.Path)
	p, _ := URL2Page(fpath)
	switch r.Method {
	case http.MethodGet:
		if err := p.MarshalEdit(w); err != nil {
			http.Error(w, "Bad Request", 400)
			log.Println("editHandler: Error in marshaling page:", err)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println("editHandler: ParseForm failed bailing:", err)
			http.Error(w, "Bad Request", 400)
			return
		}
		text, ok := r.Form["text"]
		if !ok || len(text) < 1 {
			log.Println("editHandler: malformed post bailing")
			http.Error(w, "Bad Request", 400)
			return
		}
		p.Text = text[0]
		err = p.Save()
		if err != nil {
			log.Println("editHandler: error saving page:", err)
		}
		http.Redirect(w, r, path.Join("/", fpath), http.StatusSeeOther)
	default:
		http.Error(w, "Bad Request", 400)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = frontPage
	}
	p, _ := URL2Page(path)
	if err := p.Marshal(w); err != nil {
		http.Error(w, "Bad Request", 400)
		log.Println("pageHandler: Error in marshalling page:", err)
	}
}
