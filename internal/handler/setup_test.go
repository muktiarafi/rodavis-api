package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/muktiarafi/rodavis-api/internal/driver"
	"github.com/muktiarafi/rodavis-api/internal/model"
	"github.com/muktiarafi/rodavis-api/internal/repository"
	"github.com/muktiarafi/rodavis-api/internal/service"
	"github.com/muktiarafi/rodavis-api/internal/utils"
	"github.com/muktiarafi/rodavis-api/internal/validation"
	"github.com/ory/dockertest/v3"
)

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

var router *chi.Mux

var (
	imagePath  string
	savePath   string
	assetsPath string
)

func TestMain(m *testing.M) {
	router = chi.NewRouter()
	db := &driver.DB{
		SQL: newTestDatabase(),
	}

	userRepo := repository.NewUserRepository(db)

	pwd, _ := os.Getwd()
	assetsPath = filepath.Join(pwd, "..", "..", "assets")
	imagePath = filepath.Join(assetsPath, "images")
	savePath = filepath.Join(imagePath, "tests")
	imgPersist := utils.NewLocalImagePersistence(savePath)
	userSRV := service.NewUserService(userRepo, imgPersist)

	v := validator.New()
	trans := validation.NewDefaultTranslator(v)
	val := validation.NewValidator(v, trans)
	userHandler := NewUserHandler(val, userSRV)
	userHandler.Route(router)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func newTestDatabase() *sql.DB {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err = pool.Run("postgres", "alpine", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		db, err = sql.Open(
			"pgx",
			fmt.Sprintf("host=localhost port=%s dbname=postgres user=postgres password=secret", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}

		migrationFilePath := filepath.Join("..", "..", "db", "migrations")
		return driver.Migration(migrationFilePath, db)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db
}

func register(createUserDTO *model.CreateUserDTO) (*model.UserDTO, *httptest.ResponseRecorder) {
	b, _ := json.Marshal(createUserDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	responseBody, _ := ioutil.ReadAll(res.Body)

	apiResponse := struct {
		Data *model.UserDTO `json:"data"`
	}{}
	json.Unmarshal(responseBody, &apiResponse)

	return apiResponse.Data, res
}

func assertResponseCode(t testing.TB, want, got int) {
	t.Helper()

	if got != want {
		t.Errorf("Expecting status code %d, but got %d instead", want, got)
	}
}
