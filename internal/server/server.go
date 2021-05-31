package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

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
	"github.com/muktiarafi/rodavis-api/internal/validation"
)

type App struct {
	*chi.Mux
	*sql.DB
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
	userSRV := service.NewUserService(userRepo)

	val := validator.New()
	trans := validation.NewDefaultTranslator(val)
	v := validation.NewValidator(val, trans)
	userHandler := handler.NewUserHandler(v, userSRV)
	userHandler.Route(r)

	predictAPIURL := os.Getenv("PREDICT_API_URL")
	reportRepo := repository.NewReportRepository(db)
	reportSRV := service.NewReportService(reportRepo, userRepo, predictAPIURL)
	reportHandler := handler.NewReportHandler(v, reportSRV)
	reportHandler.Route(r)

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
		Mux: r,
		DB:  db.SQL,
	}
	return app
}
