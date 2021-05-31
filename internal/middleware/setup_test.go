package middleware

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/muktiarafi/rodavis-api/internal/api"
	"github.com/muktiarafi/rodavis-api/internal/utils"
)

var router *chi.Mux

func TestMain(m *testing.M) {
	router = chi.NewMux()

	os.Setenv("JWT_KEY", "12345678")

	router.With(RequireAuth).Get("/tokens", testRequireAuthHandler)

	os.Exit(m.Run())
}

func testRequireAuthHandler(w http.ResponseWriter, r *http.Request) {
	userPayload, ok := r.Context().Value("userPayload").(*utils.UserPayload)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.NewResponse(http.StatusOK, "OK", userPayload)
}
