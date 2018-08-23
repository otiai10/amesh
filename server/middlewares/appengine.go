// +build appengine

package middlewares

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// Context ...
func Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

// HTTPClient ...
func HTTPClient(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}

// Logger ...
type Logger struct {
	ctx context.Context
}

// Log ...
func Log(ctx context.Context) Logger {
	return Logger{ctx}
}

// Debugf ...
func (logger Logger) Debugf(format string, v ...interface{}) {
	log.Debugf(logger.ctx, format, v...)
}

// Errorf ...
func (logger Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(logger.ctx, format, v...)
}
