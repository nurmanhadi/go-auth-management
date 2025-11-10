package routes

import (
	"auth-management/delivery/rest/handler"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router      *chi.Mux
	UserHandler *handler.UserHandler
}

func (r *Router) New() {
	r.Router.Route("/api", func(api chi.Router) {
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/register", r.UserHandler.UserRegister)
			auth.Post("/login", r.UserHandler.UserLogin)
		})
	})
}
