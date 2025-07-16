package middleware

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	status int
}

func (w *WrapperWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}
