package main

import (
	"fmt"
	"net/http"
	"strings"
)

// https://ndersson.me/post/capturing_status_code_in_net_http/
// https://upgear.io/blog/golang-tip-wrapping-http-response-writer-for-middleware/
type LoggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (self *LoggingResponseWriter) WriteHeader(code int) {
	self.status = code
	self.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		logger.Debugf(" In %v %v %v", r.RemoteAddr, r.Method, r.URL)

		// during
		if strings.Contains(fmt.Sprintf("%v", r.URL), "/ws/") {
			next.ServeHTTP(w, r)
			return
		}
		//.end

		// Initialize the status to 200 in case WriteHeader is not called
		lrw := LoggingResponseWriter{w, 200}

		// during
		next.ServeHTTP(&lrw, r)

		// end
		logger.Debugf("Out %v %v %v [%v]", r.RemoteAddr, r.Method, r.URL, lrw.status)
	})
}

func SetHeadersMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
