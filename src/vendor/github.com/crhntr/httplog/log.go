package httplog

import (
	"log"
	"net/http"
	"os"
	"time"
)

func JSON(outLogger, errLogger *log.Logger) func(req *http.Request, elapsed time.Duration, status int) {
	return func(req *http.Request, elapsed time.Duration, status int) {
		if status >= 500 {
			errLogger.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %d}`+"\n", req.Method, req.URL.Path, elapsed, status)
		}
		outLogger.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %d}`+"\n", req.Method, req.URL.Path, elapsed, status)
	}
}

type Func func(req *http.Request, elapsed time.Duration, status int)

type logRecord struct {
	http.ResponseWriter
	status int
}

func (r *logRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

// WriteHeader implements ResponseWriter for logRecord
func (r *logRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Wrap(f http.Handler, logFns ...Func) http.HandlerFunc {
	outLogger := log.New(os.Stdout, "", 0)
	errLogger := log.New(os.Stderr, "", 0)

	var fn Func
	if len(logFns) == 0 {
		fn = JSON(outLogger, errLogger)
	} else if len(logFns) == 1 {
		fn = logFns[0]
	} else {
		fn = func(req *http.Request, elapsed time.Duration, status int) {
			for _, lg := range logFns {
				lg(req, elapsed, status)
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		start := time.Now()
		f.ServeHTTP(record, r)

		fn(r, time.Since(start), record.status)
	}
}
