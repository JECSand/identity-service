package routing

import "net/http"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// NewWriterWrapper response writer wrapper
func NewWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{ResponseWriter: w}
}

func (rw *ResponseWriterWrapper) Status() int {
	return rw.status
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}
