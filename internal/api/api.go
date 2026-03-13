package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/devvdark0/auth-service/internal/auth"
	"github.com/devvdark0/auth-service/internal/config"
	"github.com/devvdark0/auth-service/internal/handler"
	"github.com/devvdark0/auth-service/internal/repository/postgres"
	"github.com/devvdark0/auth-service/internal/service"
	"github.com/gorilla/mux"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type API struct {
	srv *http.Server
}

func New(cfg *config.Config) *API {
	log := configureLogger(cfg.App.Env)

	log.Debug("initializing db....")
	db, err := config.InitDb(&cfg.Db)
	if err != nil {
		log.Error("failed to init databaase", "err", err)
	}

	jwtManager := auth.NewJWTManager(&cfg.Auth)
	authStorage := postgres.NewAuthRepository(db, log)
	authService := service.NewAuthService(authStorage, jwtManager)
	authHandler := handler.NewAuthHandler(log, authService)

	r := mountRoutes(log, *authHandler, jwtManager)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	return &API{
		srv: &srv,
	}
}

func (a *API) Run() error {
	if err := a.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *API) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}

func mountRoutes(log *slog.Logger, handler handler.AuthHandler, jwt *auth.JWTManager) *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware(log))

	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	authRoutes.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

	protectedRoutes := r.PathPrefix("/api/v1").Subrouter()
	protectedRoutes.Use(AuthMiddleware(jwt))
	protectedRoutes.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, I'm authorized!!")
	}).Methods(http.MethodGet)

	return r
}

func configureLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
