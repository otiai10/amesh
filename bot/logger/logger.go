package logger

import (
	"context"

	"cloud.google.com/go/logging"
)

// Client ...
type Client struct {
	Name string
	*logging.Client
}

// New ...
// User need to close this client.
func New(ctx context.Context, project string, name string) (Client, error) {
	client, err := logging.NewClient(ctx, project)
	return Client{
		Name:   name,
		Client: client,
	}, err
}

// Debug ...
func (lg Client) Debug(entry interface{}) {
	lg.Client.Logger(lg.Name).Log(logging.Entry{Payload: entry, Severity: logging.Debug})
}

// Infof ...
func (lg Client) Infof(format string, args ...interface{}) {
	lg.Client.Logger(lg.Name).StandardLogger(logging.Info).Printf(format, args...)
}

// Errorf ...
func (lg Client) Errorf(format string, args ...interface{}) {
	lg.Client.Logger(lg.Name).StandardLogger(logging.Error).Printf(format, args...)
}

// Criticalf ...
func (lg Client) Criticalf(format string, args ...interface{}) {
	lg.Client.Logger(lg.Name).StandardLogger(logging.Critical).Printf(format, args...)
}
