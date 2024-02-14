package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Info(namespace, message string) {
	log.Info().Msg(fmt.Sprintf("[%s] %s", namespace, message))
}

func Error(namespace, message string, err ...error) {
	var stringBuilder []string

	stringBuilder = append(stringBuilder, fmt.Sprintf("[%s] %s", namespace, message))

	if len(err) > 0 {
		stringBuilder = append(stringBuilder, fmt.Sprintf(": %s", err[0].Error()))
	}

	log.Error().Msg(strings.Join(stringBuilder, " "))
}
