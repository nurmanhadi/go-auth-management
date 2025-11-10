package main

import (
	"auth-management/config"
	"net/http"
)

func main() {
	config.NewEnv()
	logger := config.NewLogger()
	validator := config.NewValidator()
	db := config.NewDatabase()
	r := config.NewRouter()
	config.Initialize(&config.Bootstrap{
		DB:        db,
		Logger:    logger,
		Router:    r,
		Validator: validator,
	})

	err := http.ListenAndServe("0.0.0.0:3000", r)
	if err != nil {
		panic(err)
	}
}
