package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/middleware"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/service"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/validation"
)

type UserHandler struct {
	*validation.Validator
	service.UserService
}

func NewUserHandler(validator *validation.Validator, userSRV service.UserService) *UserHandler {
	return &UserHandler{
		Validator:   validator,
		UserService: userSRV,
	}
}

func (h *UserHandler) Route(mux *chi.Mux) {
	mux.Route("/api/users", func(r chi.Router) {
		r.Post("/register", h.CreateUser)
		r.Post("/login", h.Auth)
		r.With(middleware.RequireAuth).Get("/", h.GetUser)
	})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserDTO := new(model.CreateUserDTO)
	if err := api.Bind(r.Body, createUserDTO); err != nil {
		api.SendError(w, err)
		return
	}

	if err := h.Validate("UserHandler.CreateUser", createUserDTO); err != nil {
		api.SendError(w, err)
		return
	}

	user, err := h.UserService.Create(r.Context(), createUserDTO)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusCreated, "Created", user).SendJSON(w)
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	loginDTO := new(model.LoginDTO)
	if err := api.Bind(r.Body, loginDTO); err != nil {
		api.SendError(w, err)
		return
	}

	if err := h.Validate("UserHandler.Auth", loginDTO); err != nil {
		api.SendError(w, err)
		return
	}

	user, err := h.UserService.Auth(r.Context(), loginDTO)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusOK, "OK", user).SendJSON(w)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userPayload, err := api.UserPayloadFromContext("UserHandler.GetUser", r)
	if err != nil {
		api.SendError(w, err)
		return
	}

	user, err := h.UserService.Get(r.Context(), userPayload.ID)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusOK, "OK", user).SendJSON(w)
}
