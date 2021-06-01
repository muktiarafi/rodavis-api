package api

import (
	"context"
	"errors"
	"net/http"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

type key struct{}

func UserPayloadFromContext(op string, r *http.Request) (*model.UserPayload, error) {
	userPayload, ok := r.Context().Value(key{}).(*model.UserPayload)
	if !ok {
		return nil, &Exception{
			Op:  op,
			Err: errors.New("missing user payload"),
		}
	}

	return userPayload, nil
}

func UserPayloadToContext(userPayload *model.UserPayload, r *http.Request) context.Context {
	return context.WithValue(r.Context(), key{}, userPayload)
}
