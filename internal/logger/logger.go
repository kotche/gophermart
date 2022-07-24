package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func Init() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger
}
