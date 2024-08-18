package logger

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type LogVariable = map[string]any

type LogHandler struct {
	Instance zerolog.Logger
}

type logStack struct {
	logger LogHandler
	Writer io.Closer
}

var LogStack = map[string]logStack{}

func ClearLogStack() {
	LogStack = map[string]logStack{}
}

func Register() {
	InitAppLog()
}

func NewConsole() LogHandler {
	var stdOut = os.Stdout
	outputConsole := zerolog.ConsoleWriter{Out: stdOut, TimeFormat: time.RFC3339}
	outputConsole.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(
			fmt.Sprintf("%8s", fmt.Sprintf("[%s]:", i)),
		)
	}
	logger := zerolog.New(outputConsole).With().Timestamp().Logger()
	randomizer := rand.New(rand.NewSource(10))
	randomInt := fmt.Sprintf("%d", randomizer.Uint32())

	handler := LogHandler{logger}

	LogStack[randomInt] = logStack{
		logger: handler,
		Writer: stdOut,
	}

	return handler
}

func NewWithWriters(writers ...io.Writer) zerolog.Logger {
	multiOutput := io.MultiWriter(writers...)
	return zerolog.New(multiOutput).With().Timestamp().Logger()
}

func (l LogHandler) Warn(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Warn().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Warnf(format string, m ...any) {
	l.Instance.Warn().Msgf(format, m...)
}

func (l LogHandler) Error(e error, m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Error().Err(e).Msg(strings.Join(messages, " "))
}

func (l LogHandler) Errorf(e error, format string, m ...any) {
	l.Instance.Error().Err(e).Msgf(format, m...)
}

func (l LogHandler) Fatal(e error, m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Fatal().Err(e).Msg(strings.Join(messages, " "))
}

func (l LogHandler) Fatalf(e error, format string, m ...any) {
	l.Instance.Fatal().Err(e).Msgf(format, m...)
}

func (l LogHandler) Info(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Info().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Infof(format string, m ...any) {
	l.Instance.Info().Msgf(format, m...)
}

func (l LogHandler) Debug(m ...any) {
	var messages = l.messageToString(m)
	l.Instance.Debug().Msg(strings.Join(messages, " "))
}

func (l LogHandler) Debugf(format string, m ...any) {
	l.Instance.Debug().Msgf(format, m...)
}

func (l LogHandler) InfoWithVariables(vars LogVariable, m ...any) {
	var messages = l.messageToString(m)
	logEvent := l.Instance
	for key, value := range vars {
		variableValue := fmt.Sprintf("%v", value)
		logEvent = logEvent.With().Str(key, variableValue).Logger()
	}
	logEvent.Info().Msg(strings.Join(messages, " "))
}

func (l LogHandler) DebugWithVariables(vars LogVariable, m ...any) {
	var messages = l.messageToString(m)
	logEvent := l.Instance
	for key, value := range vars {
		variableValue := fmt.Sprintf("%v", value)
		logEvent = logEvent.With().Str(key, variableValue).Logger()
	}
	logEvent.Debug().Msg(strings.Join(messages, " "))
}

func (l LogHandler) messageToString(m []any) []string {
	var messages []string
	for _, msg := range m {
		messages = append(messages, fmt.Sprintf("%v", msg))
	}
	return messages
}
