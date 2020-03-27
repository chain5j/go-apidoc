package middleware

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
)

//go:generate easytags $GOFILE

type ResponseRecorder struct {
	writer http.ResponseWriter
	Status int
	Body   *bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	r := &ResponseRecorder{
		writer:w,
		Status:http.StatusOK,
		Body:bytes.NewBuffer(nil),
	}
	return r
}

func (r *ResponseRecorder) Header() http.Header {
	return r.writer.Header()
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.writer.WriteHeader(status)
}

func (r *ResponseRecorder) Write(buf []byte) (int, error) {
	n, err := r.writer.Write(buf)
	if err == nil {
		r.Body.Write(buf)
	}
	return n, err
}

func (r *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.writer.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("Error in hijacker")
}