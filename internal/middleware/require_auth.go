package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/muktiarafi/rodavis-api/internal/api"
	"github.com/muktiarafi/rodavis-api/internal/utils"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			exc := api.NewSingleMessageException(
				api.EUNAUTHORIZED,
				"RequiredAuth",
				"Not Authorized",
				errors.New("missing value in authorization header"),
			)
			api.SendError(w, exc)
			return
		}

		token, payload, err := utils.ParseToken(authHeader[7:])
		if err != nil {
			exc := api.NewSingleMessageException(
				api.EUNAUTHORIZED,
				"RequiredAuth",
				"Not Authorized",
				errors.New("invalid token"),
			)
			api.SendError(w, exc)
			return
		}

		if !token.Valid {
			exc := api.NewSingleMessageException(
				api.EUNAUTHORIZED,
				"RequiredAuth",
				"Not Authorized",
				errors.New("invalid token"),
			)
			api.SendError(w, exc)
			return
		}

		ctx := api.UserPayloadToContext(payload, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
