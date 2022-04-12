package main

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/hello-release/src/cmd/hello-server/fakes"
)

func Test_requireGetMethod(t *testing.T) {
	var nextCallCount int
	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
		nextCallCount++
	})

	h := requireGetMethod(next)
	resetCallCount := func() { nextCallCount = 0 }

	t.Run(http.MethodGet, func(t *testing.T) {
		t.Cleanup(resetCallCount)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if nextCallCount != 1 {
			t.Error("expected next to be called")
		}

		res := rec.Result()
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("expected next's status code, got: %d", res.StatusCode)
		}
	})

	t.Run("not GET", func(t *testing.T) {
		t.Cleanup(resetCallCount)

		req, _ := http.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if nextCallCount != 0 {
			t.Error("expected next to not be called")
		}

		res := rec.Result()
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected method not allowed, got: %d", res.StatusCode)
		}
	})
}

func Test_indexHTML(t *testing.T) {
	_, err := template.New("index.html").Parse(indexPageSource)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func Test_indexPage(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		templates := new(fakes.Execute)
		templates.ExecuteReturns(nil)
		templates.ExecuteStub = func(w io.Writer, d interface{}) error {
			_, _ = w.Write([]byte("output"))
			return nil
		}
		var logBuf bytes.Buffer
		logger := log.New(&logBuf, "", 0)

		h := indexPage(templates, logger)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got: %d", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
		}
		if !bytes.Equal(body, []byte("output")) {
			t.Errorf("unexpected body, got: %q", string(body))
		}
		if logBuf.Len() != 0 {
			t.Errorf("unexpected no logs, got: %q", logBuf.String())
		}
	})

	t.Run("successful render", func(t *testing.T) {
		templates := new(fakes.Execute)
		templates.ExecuteReturns(errors.New("banana"))

		var logBuf bytes.Buffer
		logger := log.New(&logBuf, "", 0)

		h := indexPage(templates, logger)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		res := rec.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status %d, got: %d", http.StatusInternalServerError, res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
		}
		if bytes.Contains(body, []byte("banana")) {
			t.Errorf("did not expect internal template error, got: %q", string(body))
		}
		if !bytes.Contains(body, []byte("failed to render template")) {
			t.Errorf("expected semi helpful error, got: %q", string(body))
		}

		logs, err := io.ReadAll(&logBuf)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
		}
		if !bytes.Contains(logs, []byte("banana")) {
			t.Errorf("expected template error, got: %q", string(body))
		}
	})
}
