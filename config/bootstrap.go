package config

import (
	"auth-management/delivery/rest/handler"
	"auth-management/delivery/rest/middleware"
	"auth-management/delivery/rest/routes"
	"auth-management/internal/cache"
	"auth-management/internal/repository"
	"auth-management/internal/service"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Bootstrap struct {
	DB        *gorm.DB
	Cache     *memcache.Client
	Logger    zerolog.Logger
	Validator *validator.Validate
	Router    *chi.Mux
}

func Initialize(deps *Bootstrap) {
	// cache
	tokeCache := cache.NewTokenCache(deps.Cache)

	// repository
	userRepo := repository.NewUserRepository(deps.DB)

	// service
	userServ := service.NewUserService(deps.Logger, deps.Validator, userRepo, tokeCache)

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
