package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/tonievictor/dotenv"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/database"
	"ogugu/docs"
	"ogugu/models"
	"ogugu/router"
	"ogugu/telemetry"
)

func main() {
	dotenv.Config()

	docs.SwaggerInfo.Title = "Ogugu"
	docs.SwaggerInfo.Description = "An RSS feed reader"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "localhost:8080" // this will be dynamic
	docs.SwaggerInfo.BasePath = "/v1/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	logger, _ := zap.NewProduction()
	db, err := database.Init("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("unable to initialize database")
		return
	}

	models.New(db)

	InitServer(logger)
}

func InitServer(log *zap.Logger) {
	r := router.Routes(log)

	log.Info("Server is running on port 8080")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// set up tracing
	exp, err := telemetry.NewConsoleExporter()
	tp := telemetry.NewTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err = <-serverErr:
		log.Error("API", zap.String("error", err.Error()))
		close(serverErr)
		return
	case <-ctx.Done():
		log.Info("API: Graceful shutdown: received\n")
		stop()
	}

	server.Shutdown(ctx)
	return
}
