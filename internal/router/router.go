package router

import (
	"io"
	"net/http"

	"github.com/devvdark0/auth-service/internal/auth"
	handler "github.com/devvdark0/auth-service/internal/handler/auth"
	"github.com/devvdark0/auth-service/internal/logger"
	"github.com/devvdark0/auth-service/internal/middleware"
	"github.com/gorilla/mux"
)

func InitRoutes(handler handler.AuthHandler, log *logger.Logger, manager *auth.JWTManager) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware(log))

	authRouter := r.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", handler.Register)
	authRouter.HandleFunc("/login", handler.Login)

	userRouter := r.PathPrefix("/api/v1").Subrouter()
	userRouter.Use(middleware.AuthMiddleware(manager))
	userRouter.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, I'm authorized!!")
	})

	return r
}
