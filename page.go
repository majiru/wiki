package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
    	"github.com/gomarkdown/markdown"
)

var pageFinder *regexp.Regexp

func init(){
	pageFinder = regexp.MustCompile(`([A-Z][a-z]+)+`)
}

func listPages() ([]string, error) {
	var out []string
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		out = append(out, f.Name())
	}
	return out, nil
}

func MarshalSearch(w io.Writer, query string) error {
	files, err := listPages()
	if err != nil {
		return err
	}
	list := struct { Pages []string }{}
	for _, f := range files {
		if strings.Contains(f, query) {
			list.Pages = append(list.Pages, f)
		}
	}
	t, err := template.New("search").Parse(searchPage)
	if err != nil {
		return err
	}
	return t.Execute(w, list)
}

type Page struct {
	Text, Title, Body string
}

func URL2Page(path string) (*Page, error) {
	fpath := filepath.Join(dataDir, path)
	p := &Page{Title: filepath.Base(fpath)}
	f, err := os.Open(fpath)
	if err != nil {
		return p, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return p, err
	}
	p.Text = string(content)
	return p, nil
}

func (p *Page) Save() error {
	f, err := os.Create(filepath.Join(dataDir, p.Title))
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(p.Text))
	return err
}

func (p *Page) Marshal(w io.Writer) error {
	t, err  := template.New("page").Parse(page)
	if err != nil {
		return err
	}
	b := pageFinder.ReplaceAllString(p.Text, `<a href="/$0">$0</a>`)
	p.Body = string(bluemonday.UGCPolicy().SanitizeBytes(markdown.ToHTML([]byte(b), nil, nil)))
	err = t.Execute(w, p)
	if err != nil {
		return err
	}
	return err
}

func (p *Page) MarshalEdit(w io.Writer) error {
	t, err  := template.New("edit").Parse(editPage)
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}

const page = `
<html>
<head>
	<title>{{.Title}}</title>
</head>
<body>
	<hr>
	{{- .Body -}}
	<hr>
	<a href="/edit/{{.Title}}">Edit</a>
	<form action="/search" method="POST">
	<input type="text" name="query">
	<input type="submit" value="search">
	</form>
</body>
</html>
`

const editPage = `
<html>
<head>
	<title>{{.Title}}</title>
</head>
<body>
	<form action="/edit/{{.Title}}" method="POST">
		<textarea name="text" rows=25 cols=80>{{.Text}}</textarea><br>
		<input type="submit" value="save">
	</form>
</body>
</html>
`

const searchPage=`
<html>
<head>
	<title>Search</title>
</head>
<body>
{{range .Pages}}<a href="/{{.}}">{{.}}</a>{{end}}
</body>
</html>
`
