// @Title Subscription Service API
// @Description REST API for managing user subscriptions
// @Version 1.0
// @Host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
	"subscription-service/internal/config"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер pgx для sqlx
	_ "subscription-service/docs"      // swagger docs
)

func main() {
	
	cfg := config.Load()

	
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		logger.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	svc := service.NewSubscriptionService(repo)

	subHandler := handlers.NewSubscriptionHandler(svc, logger)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", subHandler.Create)
		r.Get("/", subHandler.List)
		r.Get("/total", subHandler.SumByPeriod)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", subHandler.Get)
			r.Put("/", subHandler.Update)
			r.Delete("/", subHandler.Delete)
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced shutdown", "error", err)
	}
	logger.Info("Server stopped")
}