package main

import (
	"backend/pkg/config"
	"backend/pkg/product"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// postgresql://admin:%v@localhost:5438/Coursework for local use
	// postgresql://admin:%v@postgres:5432/Coursework for container deploy
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgresql://admin:%v@postgres:5432/Coursework", cfg.DatabasePass))
	if err != nil {
		log.Error("unable to connect to postgres")
		return
	}

	// localhost:6379 for local use
	// redis:6379 for container deploy
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Error("unable to connect to redis")
		return
	}

	productRepo := product.NewRepository(pool, client, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(cfg.Http.Timeout))
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/product", product.NewFindAll(log, productRepo))
	router.Get("/product/{productId}", product.NewFindOne(log, productRepo))
	router.Post("/product", product.NewCreate(log, productRepo))
	router.Put("/product/{productId}", product.NewUpdate(log, productRepo))
	router.Delete("/product/{productId}", product.NewDelete(log, productRepo))

	srv := &http.Server{
		Addr:         cfg.Http.Address,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		IdleTimeout:  cfg.Http.Timeout,
		Handler:      router,
	}

	log.Info("Starting listening on " + srv.Addr)
	if err = srv.ListenAndServe(); err != nil {
		log.Error("can't open server", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
