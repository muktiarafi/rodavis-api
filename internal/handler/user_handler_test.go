package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

func TestUserHandlerCreate(t *testing.T) {
	t.Run("Create user with valid payload", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambank",
			PhoneNumber: "+6212345678910",
			Email:       "bambank@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		if userDTO.Token == "" {
			t.Error("Expecting token but get none")
		}

		if userDTO.User.Email != createUserDTO.Email {
			t.Errorf("Expecting email to be %q, but got %q instead", userDTO.User.Email, createUserDTO.Email)
		}
	})

	t.Run("Create User with invalid name", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "b",
			Email:       "bambank@gmail.com",
			PhoneNumber: "+6212345678977",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Create User with invalid email", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bwerwer",
			Email:       "bambankgmailcom",
			PhoneNumber: "+6212345348910",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Create User with invalid password", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bwerwer",
			Email:       "bambankgmailcom",
			PhoneNumber: "+621234567899",
			Password:    "123",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Create User with invalid phonenumber", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bwerwer",
			Email:       "bambankgmailcom",
			PhoneNumber: "+6212",
			Password:    "123",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Create User with invalid payload", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "",
			Email:       "",
			Password:    "",
			PhoneNumber: "",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Create User with duplicate email", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "paijo",
			Email:       "paijo@gmail.com",
			PhoneNumber: "+6213445678910",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		_, res = register(createUserDTO)

		assertResponseCode(t, http.StatusConflict, res.Code)
	})

	t.Run("Create User with duplicate phonenumber", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "lala",
			Email:       "popopqewmwqewqe@gmail.com",
			PhoneNumber: "+6213445678916",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		_, res = register(createUserDTO)

		assertResponseCode(t, http.StatusConflict, res.Code)
	})
}

func TestUserHandlerAuth(t *testing.T) {
	t.Run("Login with valid credential", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "brando",
			Email:       "brando@gmail.com",
			PhoneNumber: "+6212334678910",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		loginDTO := &model.LoginDTO{
			Email:    createUserDTO.Email,
			Password: createUserDTO.Password,
		}
		b, _ := json.Marshal(loginDTO)
		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusOK, res.Code)

		responseBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data *model.UserDTO `json:"data"`
		}{}

		json.Unmarshal(responseBody, &apiResponse)

		if apiResponse.Data.Token == "" {
			t.Error("Expecting token but get none")
		}

		if apiResponse.Data.User.Email != loginDTO.Email {
			t.Errorf("Expecting email to be %q, but got %q instead", apiResponse.Data.User.Email, loginDTO.Email)
		}
	})

	t.Run("Login with invalid email", func(t *testing.T) {
		loginDTO := &model.LoginDTO{
			Email:    "brandogmailcom",
			Password: "12345678",
		}
		b, _ := json.Marshal(loginDTO)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		loginDTO := &model.LoginDTO{
			Email:    "brando@gmail.com",
			Password: "",
		}
		b, _ := json.Marshal(loginDTO)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Login with invalid payload", func(t *testing.T) {
		loginDTO := &model.LoginDTO{
			Email:    "",
			Password: "",
		}
		b, _ := json.Marshal(loginDTO)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Login with email that not been used", func(t *testing.T) {
		loginDTO := &model.LoginDTO{
			Email:    "werwerdsfewrwe@gmail.com",
			Password: "werwerwerwerwerwer",
		}
		b, _ := json.Marshal(loginDTO)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "yaya",
			Email:       "yaya@gmail.com",
			PhoneNumber: "+62127734678910",
			Password:    "12345678",
		}
		_, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		loginDTO := &model.LoginDTO{
			Email:    createUserDTO.Email,
			Password: "boboboy",
		}
		b, _ := json.Marshal(loginDTO)
		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})
}

func TestUserHandlerGetUser(t *testing.T) {
	t.Run("Get user that just created", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "werwerwet",
			Email:       "lololololo@gmail.com",
			PhoneNumber: "+6212334678912",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusOK, res.Code)

		responseBody, _ := ioutil.ReadAll(res.Body)

		getUserResponseBody := struct {
			Data *entity.User `json:"data"`
		}{}

		json.Unmarshal(responseBody, &getUserResponseBody)

		got := getUserResponseBody.Data.ID
		want := userDTO.User.ID

		if got != want {
			t.Errorf("Expecting id to be %d, but got %d instead", want, got)
		}
	})
}
