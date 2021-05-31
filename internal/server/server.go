package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/muktiarafi/rodavis-api/internal/api"
	"github.com/muktiarafi/rodavis-api/internal/config"
	"github.com/muktiarafi/rodavis-api/internal/driver"
	"github.com/muktiarafi/rodavis-api/internal/handler"
	"github.com/muktiarafi/rodavis-api/internal/logger"
	"github.com/muktiarafi/rodavis-api/internal/middleware"
	"github.com/muktiarafi/rodavis-api/internal/model"
	"github.com/muktiarafi/rodavis-api/internal/repository"
	"github.com/muktiarafi/rodavis-api/internal/service"
	"github.com/muktiarafi/rodavis-api/internal/utils"
	"github.com/muktiarafi/rodavis-api/internal/validation"
	"google.golang.org/api/option"
)

type App struct {
	*chi.Mux
	*sql.DB
	*storage.Client
}

func New() *App {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)

	const op = "server.New"
	logger.Notice(op, "Connecting to Database")
	db, err := driver.ConnectSQL(config.PostgresDSN(), true)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		logger.Error(op, &model.SourceLocation{
			File:     file,
			Function: "driver.ConnectSQL",
			Line:     line,
		}, err)
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(db)

	ctx := context.Background()
	logger.Notice(op, "Connecting to Cloud Storage")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("key.json"))
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		logger.Error(op, &model.SourceLocation{
			File:     file,
			Line:     line,
			Function: "storage.NewClient",
		}, err)
		log.Fatal(err)
	}
	bucketName := os.Getenv("BUCKET_NAME")
	bucket := client.Bucket(bucketName)
	cloudImgPersistence := utils.NewCloudImagePersistence(ctx, bucket, bucketName)
	userSRV := service.NewUserService(userRepo, cloudImgPersistence)

	val := validator.New()
	trans := validation.NewDefaultTranslator(val)
	v := validation.NewValidator(val, trans)
	userHandler := handler.NewUserHandler(v, userSRV)
	userHandler.Route(r)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := api.NewSingleMessageException(
			api.ENOTFOUND,
			"r.NotFound",
			"Route Not Found",
			errors.New("route not found"),
		)
		api.NewErrorResponse(err).SendJSON(w)
	})

	app := &App{
		Mux:    r,
		DB:     db.SQL,
		Client: client,
	}
	return app
}
