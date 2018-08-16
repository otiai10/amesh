// +build !appengine

package middlewares

import (
	"context"
	"log"
	"net/http"
)

// Context ...
func Context(r *http.Request) context.Context {
	return context.Background()
}

// HTTPClient ...
func HTTPClient(ctx context.Context) *http.Client {
	return http.DefaultClient
}

// Logger ...
type Logger struct {
}

// Log ...
func Log(ctx context.Context) Logger {
	return Logger{}
}

// Debugf ...
func (logger Logger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
