package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muktiarafi/rodavis-api/internal/model"
	"github.com/muktiarafi/rodavis-api/internal/utils"
)

func TestRequireAuth(t *testing.T) {
	t.Run("Pass valid token", func(t *testing.T) {
		userPayload := &model.UserPayload{
			ID:    1,
			Email: "bambank@gmai.com",
		}
		token, _ := utils.CreateToken(userPayload)

		req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("Expecting status code to be %d, but got %d instead", want, got)
		}
	})

	t.Run("Request without Authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusUnauthorized

		if got != want {
			t.Errorf("Expecting status code to be %d, but got %d instead", want, got)
		}
	})

	t.Run("Pass invalid token", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIsImVtYWlsIjoiYXlheWF5YUBnbWFpLmNvbSIsImV4cCI6MTYyMzE1ODg4MiwiaWF0IjoxNjIxOTQ5MjgyLCJzdWIiOiIxIn0.xLoEEnROiEPidGQJZuwigv1fthJ_eNvNr9fBVwHA1u8"
		req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		got := res.Code
		want := http.StatusUnauthorized

		if got != want {
			t.Errorf("Expecting status code to be %d, but got %d instead", want, got)
		}
	})
}
