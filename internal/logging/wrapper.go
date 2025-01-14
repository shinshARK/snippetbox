package logging

import (
	"log/slog"
)

// slogWrapper wraps slog.Logger for compatibility with log.Logger.
type slogWrapper struct {
	logger *slog.Logger
}

// Write implements the io.Writer interface to redirect log.Logger output to slog.
func (w *slogWrapper) Write(p []byte) (n int, err error) {
	w.logger.Error(string(p)) // Redirect log output to slog at error level.
	return len(p), nil
}

// NewSlogWrapper creates a new slogWrapper instance.
func NewSlogWrapper(logger *slog.Logger) *slogWrapper {
	return &slogWrapper{logger: logger}
}
