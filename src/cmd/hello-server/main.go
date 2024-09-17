package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	routes(mux)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		execute(res, templates.Lookup("index.html"), http.StatusOK, struct{}{})
	})
}

var (
	//go:embed index.gohtml
	indexPageSource string

	templates = template.Must(template.New("index.html").Parse(indexPageSource))
)

func execute(res http.ResponseWriter, t *template.Template, code int, data any) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(code)
	_, _ = buf.WriteTo(res)
}
