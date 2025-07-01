package middlewares

import "net/http"

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write implements the http.ResponseWriter interface, recording the size of the written response.
// It delegates the write operation to the underlying ResponseWriter and tracks the number of bytes written.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader implements the http.ResponseWriter interface, recording the HTTP status code.
// It delegates setting the status code to the underlying ResponseWriter and stores the status code.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
