package logger

import (
	"github.com/muktiarafi/rodavis-api/internal/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type HTTPRequest struct {
	RequestMethod string
	RequestURL    string
	Status        int
	Latency       string
	ResponseSize  int64
}

func init() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"
	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelWarnValue = "WARNING"
	zerolog.LevelErrorValue = "ERROR"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.ErrorFieldName = "message"
}

func NewInfo() *zerolog.Event {
	return log.Info()
}

func Request(httpRequest *HTTPRequest) {
	log.Log().
		Dict("httpRequest",
			zerolog.Dict().
				Str("requestMethod", httpRequest.RequestMethod).
				Str("requestUrl", httpRequest.RequestURL).
				Int("status", httpRequest.Status).
				Str("latency", httpRequest.Latency).
				Int64("responseSize", httpRequest.ResponseSize)).
		Send()
}

func Notice(op, message string) {
	log.Log().
		Str("severity", "NOTICE").
		Dict("logging.googleapis.com/operation", zerolog.Dict().
			Str("id", op)).
		Msg(message)

}

func Info(op, message string) {
	log.Info().
		Dict("logging.googleapis.com/operation", zerolog.Dict().
			Str("id", op)).
		Msg(message)
}

func NewWarn() *zerolog.Event {
	return log.Warn()
}

func NewError(err error) *zerolog.Event {
	return log.Err(err)
}

func Error(op string, location *model.SourceLocation, err error) {
	log.Err(err).
		Dict("logging.googleapis.com/operation", zerolog.Dict().
			Str("id", op)).
		Dict("sourceLocation",
			zerolog.Dict().
				Str("file", location.File).
				Str("function", location.Function).
				Int("line", location.Line),
		).
		Send()
}
