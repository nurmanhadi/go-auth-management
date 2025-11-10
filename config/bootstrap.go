package config

import (
	"auth-management/delivery/rest/handler"
	"auth-management/delivery/rest/middleware"
	"auth-management/delivery/rest/routes"
	"auth-management/internal/repository"
	"auth-management/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Bootstrap struct {
	DB        *gorm.DB
	Logger    zerolog.Logger
	Validator *validator.Validate
	Router    *chi.Mux
}

func Initialize(deps *Bootstrap) {
	// repository
	userRepo := repository.NewUserRepository(deps.DB)

	// service
	userServ := service.NewUserService(deps.Logger, deps.Validator, userRepo)

	// handler
	userHand := handler.NewUserHandler(userServ)

	// middleware
	deps.Router.Use(middleware.ErrorHandler)

	// routes
	r := routes.Router{
		Router:      deps.Router,
		UserHandler: userHand,
	}
	r.New()
}
