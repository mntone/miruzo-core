package middleware

import (
	"log"
	"net/http"
)

type requestLogger interface {
	Printf(format string, v ...any)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func newLoggingResponseWriter(responseWriter http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: responseWriter,
		statusCode:     http.StatusOK,
	}
}

func (writer *loggingResponseWriter) WriteHeader(statusCode int) {
	if writer.wroteHeader {
		return
	}

	writer.statusCode = statusCode
	writer.wroteHeader = true
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *loggingResponseWriter) Write(body []byte) (int, error) {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}

	return writer.ResponseWriter.Write(body)
}

func logRequest(
	logger requestLogger,
	request *http.Request,
	statusCode int,
) {
	remoteAddress := request.RemoteAddr
	if remoteAddress == "" {
		remoteAddress = "-"
	}

	path := request.URL.Path
	if path == "" {
		path = "/"
	}

	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "Unknown Status"
	}

	logger.Printf(
		"[miruzo-api] %s - \"%s %s\" %d %s",
		remoteAddress,
		request.Method,
		path,
		statusCode,
		statusText,
	)
}

func RequestLogWithLogger(next http.Handler, logger requestLogger) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		loggedResponseWriter := newLoggingResponseWriter(responseWriter)
		next.ServeHTTP(loggedResponseWriter, request)

		logRequest(logger, request, loggedResponseWriter.statusCode)
	})
}

func RequestLog(next http.Handler) http.Handler {
	return RequestLogWithLogger(next, log.Default())
}
