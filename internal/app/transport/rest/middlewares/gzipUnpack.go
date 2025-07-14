package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read implements the io.Reader interface for the compressReader, reading decompressed data from the gzip reader.
// It reads into the provided byte slice and returns the number of bytes read and any error encountered.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes both the underlying reader and the gzip reader, ensuring proper resource cleanup.
// It returns an error if closing either reader fails.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func (m *middleware) gzipUnpack(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Encoding") != "gzip" {
				next.ServeHTTP(w, r)
				return
			}

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()

			// передаём управление хендлеру
			next.ServeHTTP(w, r)
		})
}
