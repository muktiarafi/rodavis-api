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
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ory/dockertest/v3"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/config"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/repository"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/service"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/validation"
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

var admin *model.LoginDTO

func TestMain(m *testing.M) {
	router = chi.NewRouter()
	db := newTestDatabase()
	mockPredictServer := mockPredictServer()

	configApp := &config.App{
		DB: db,
	}
	userRepo := repository.NewUserRepository()

	pwd, _ := os.Getwd()
	assetsPath = filepath.Join(pwd, "..", "..", "assets")
	imagePath = filepath.Join(assetsPath, "images")
	savePath = filepath.Join(imagePath, "tests")
	userSRV := service.NewUserService(configApp, userRepo)

	v := validator.New()
	trans := validation.NewDefaultTranslator(v)
	val := validation.NewValidator(v, trans)
	userHandler := NewUserHandler(val, userSRV)
	userHandler.Route(router)

	reportRepo := repository.NewReportRepository()
	predictAPIURL := mockPredictServer.URL
	reportSRV := service.NewReportService(configApp, reportRepo, userRepo, predictAPIURL)
	reportHandler := NewReportHandler(val, reportSRV)
	reportHandler.Route(router)

	adminCreateUserDTO := &model.CreateUserDTO{
		Name:        "yahahaha",
		Email:       "telolet@gmail.com",
		PhoneNumber: "+6234564876902",
		Password:    "$2y$12$.EpYweFep8NzqNz0CiVKkOOQh/MCByE7DUIIyeFo5RVs7AuYibFOu",
	}
	createdAdmin, err := createAdmin(adminCreateUserDTO, db)
	if err != nil {
		panic(err)
	}

	admin = &model.LoginDTO{
		Email:    createdAdmin.Email,
		Password: "12345678",
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	mockPredictServer.Close()
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

func mockPredictServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ada yang masuk nih")
		image, header, err := r.FormFile("image")
		if err != nil {
			api.NewResponse(http.StatusBadRequest, api.EINTERNAL, "Server Error").SendJSON(w)
			return
		}
		defer image.Close()

		filenameSplit := strings.Split(header.Filename, ".")
		format := filenameSplit[len(filenameSplit)-1]
		switch format {
		case "jpg", "jpeg", "png":
			classes := []string{"D00", "D01", "D10", "D11", "D20", "D40", "D43", "D44", "D50"}
			predictResult := &model.PredictResult{
				ImageUrl: "https://storage.googleapis.com/test/predict.jpg",
				Classes:  classes,
				Score:    90,
			}

			api.NewResponse(http.StatusOK, "OK", predictResult).SendJSON(w)
			return
		}

		api.NewResponse(http.StatusBadRequest, api.EINVALID, "Mimetype not supported").SendJSON(w)
	}))
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

func createAdmin(createUserDTO *model.CreateUserDTO, db *sql.DB) (*entity.User, error) {

	stmt := `INSERT INTO users (name, phone_number, email, password, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *`

	newUser := new(entity.User)
	if err := db.QueryRow(
		stmt,
		createUserDTO.Name,
		createUserDTO.PhoneNumber,
		createUserDTO.Email,
		createUserDTO.Password,
		"ADMIN",
	).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.PhoneNumber,
		&newUser.Email,
		&newUser.Password,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return newUser, nil
}
