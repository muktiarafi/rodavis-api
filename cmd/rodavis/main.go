package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/logger"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/server"
)

func main() {
	app := server.New()

	addr := ":8080"
	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      app.Mux,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	const op = "main"
	go func() {
		<-quit
		logger.Notice(op, "Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer app.DB.Close()

		httpServer.Shutdown(ctx)
		close(done)
	}()

	logger.Notice(op, "listening on port "+addr)
	if err := httpServer.ListenAndServe(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		logger.Error(op, &model.SourceLocation{
			File:     file,
			Function: "httpServer.ListenAndServe",
			Line:     line,
		}, err)
	}

	<-done
	logger.Notice(op, "Server Stopped")
}
