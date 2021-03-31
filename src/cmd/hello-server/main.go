package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	mux := new(http.ServeMux)

	mux.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(":" + port, requireGetMethod(mux)))
}

//go:embed index.gohtml
var indexPageSource string

func indexPage(res http.ResponseWriter, req *http.Request) {
	t, err := template.New("index.html").Parse(indexPageSource)
	if err != nil {
		log.Println("index page failed to parse:", err)
		http.Error(res, "internal error: failed to parse template", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "index.html", struct{}{})
	if err != nil {
		log.Println("index page failed to render:", err)
		http.Error(res, "internal error: failed to render template", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, _ = io.Copy(res, &buf)
}

func requireGetMethod(next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req )
	}
}