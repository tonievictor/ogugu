package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/redis/go-redis/v9"
	"github.com/tonievictor/dotenv"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"ogugu/job"
	"ogugu/database"
	"ogugu/database/cache"
	_ "ogugu/docs"
	"ogugu/router"
	"ogugu/telemetry"
)

// @title Ogugu API
// @description An API for an RSS aggregator
// @version 0.1

// @host localhost:8080
// @BasePath  /v1/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your auth token in the format **Bearer &lt;token&gt;**
func main() {
	dotenv.Config()

	log, _ := zap.NewProduction()
	db, err := database.New("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("unable to initialize database", zap.String("error", err.Error()))
		return
	}
	rds, err := cache.Setup(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Error("unable to initialize cache store")
		return
	}

	go job.Start(db, log)
	InitServer(db, rds, log)
}

func InitServer(db *sql.DB, rds *redis.Client, log *zap.Logger) {
	r := router.Routes(db, rds, log)

	log.Info("Server is running on port 8080")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	exp, err := telemetry.NewOtlpExporter(ctx)
	if err != nil {
		log.Error("cannot setup otlp exporter", zap.Error(err))
		return
	}
	res, err := telemetry.NewResource()
	if err != nil {
		log.Error("cannot create telemetry resource", zap.Error(err))
		return
	}
	tp := telemetry.NewTraceProvider(res, exp)

	defer func() {
		_ = tp.Shutdown(ctx)
	}()

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
	case err := <-serverErr:
		log.Error("API", zap.String("error", err.Error()))
		close(serverErr)
		return
	case <-ctx.Done():
		log.Info("API: Graceful shutdown: received\n")
		stop()
	}

	server.Shutdown(ctx)
}
