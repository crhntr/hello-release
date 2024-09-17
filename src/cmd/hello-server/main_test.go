package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_indexPage(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		mux := http.NewServeMux()
		routes(mux)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got: %d", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
		}
		if !bytes.Contains(body, []byte("Hello, 🌎!")) {
			t.Errorf("unexpected body, got: %q", string(body))
		}
	})
}
