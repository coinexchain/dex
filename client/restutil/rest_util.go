package restutil

import (
	"net/http"
)

type ResponseWriter4UT struct {
	header     http.Header
	statusCode int
	body       []byte
}

var _ http.ResponseWriter = &ResponseWriter4UT{}

func NewResponseWriter4UT() *ResponseWriter4UT {
	return &ResponseWriter4UT{
		header:     http.Header(make(map[string][]string)),
		statusCode: 0,
		body:       make([]byte, 0, 1000),
	}
}

func (w *ResponseWriter4UT) ClearBody() {
	w.body = w.body[:0]
}

func (w *ResponseWriter4UT) GetBody() []byte {
	return w.body
}

func (w *ResponseWriter4UT) Header() http.Header {
	return w.header
}

func (w *ResponseWriter4UT) Write(bz []byte) (int, error) {
	w.body = append(w.body, bz...)
	return len(w.body), nil
}

func (w *ResponseWriter4UT) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
