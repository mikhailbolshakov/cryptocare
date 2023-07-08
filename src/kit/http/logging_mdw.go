package http

import (
	"bytes"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"io/ioutil"
	"net/http"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
			s.logger().C(r.Context()).F(log.FF{"method": r.Method, "URL": r.URL.Path, "headers": r.Header, "body": string(body)}).Trc("request")
		}
		loggableRsp := &loggableResponseWriter{ResponseWriter: w}

		next.ServeHTTP(loggableRsp, r)

		s.logger().C(r.Context()).F(log.FF{"status": loggableRsp.StatusCode, "headers": loggableRsp.Header(), "body": string(loggableRsp.Body)}).Trc("response")
	})
}

type loggableResponseWriter struct {
	http.ResponseWriter
	Body        []byte
	StatusCode  int
	wroteHeader bool
}

func (rw *loggableResponseWriter) Write(data []byte) (int, error) {
	rw.Body = append(rw.Body, data...)
	return rw.ResponseWriter.Write(data)
}

func (rw *loggableResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}
