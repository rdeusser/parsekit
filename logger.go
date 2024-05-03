package parsekit

import (
	"io"

	"github.com/rdeusser/parsekit/internal/logging"
)

var DefaultLogger = logging.New(io.Discard, logging.Info)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}
