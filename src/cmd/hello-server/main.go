package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/crhntr/httplog"
)

func main() {
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	routes(mux)

	log.Fatal(http.ListenAndServe(":"+port, httplog.Wrap(mux)))
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", indexPage(
		template.Must(template.New("index.html").Parse(indexPageSource)),
		log.Default(),
	))
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o ./fakes/execute.go --fake-name Execute . executer
//counterfeiter:generate -o ./fakes/linePrinter.go --fake-name PrintLn . linePrinter

type executer interface {
	Execute(w io.Writer, d interface{}) error
}

type linePrinter interface {
	Println(...interface{})
}

//go:embed index.gohtml
var indexPageSource string

func indexPage(t executer, logger linePrinter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		err := t.Execute(&buf, struct{}{})
		if err != nil {
			logger.Println("index page failed to render:", err)
			http.Error(res, "failed to render template", http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, _ = buf.WriteTo(res)
	}
}
