package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/config"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/handler"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/logger"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/middleware"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/repository"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/service"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/validation"
)

type App struct {
	*chi.Mux
	*sql.DB
}

func New() *App {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(mid.Recoverer)

	const op = "server.New"
	logger.Notice(op, "Connecting to Database")
	db, err := driver.ConnectSQL(config.CloudSQLConnection(), false)
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
