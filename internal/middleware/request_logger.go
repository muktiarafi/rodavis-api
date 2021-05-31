package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/muktiarafi/rodavis-api/internal/logger"
	"github.com/rs/zerolog"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)

		duration := m.Duration / time.Millisecond * time.Millisecond
		req := &logger.HTTPRequest{
			RequestMethod: r.Method,
			RequestURL:    r.RequestURI,
			Status:        m.Code,
			Latency:       duration.String(),
			ResponseSize:  m.Written,
		}

		logger.Request(req)

		if duration > (300 * time.Millisecond) {
			logger.NewWarn().
				Dict("httpRequest",
					zerolog.Dict().
						Str("requestMethod", req.RequestMethod).
						Str("requestUrl", req.RequestURL).
						Int("status", req.Status).
						Str("latency", req.Latency)).
				Msg(fmt.Sprintf("Request taking too long to complete. Duration to complete is %s", duration.String()))
		}
	})
}
