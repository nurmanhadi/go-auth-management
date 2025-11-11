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
	cache := config.NewCache()
	r := config.NewRouter()
	conn, ch := config.NewAmqp()
	defer conn.Close()
	defer ch.Close()
	config.Initialize(&config.Bootstrap{
		DB:        db,
		Cache:     cache,
		Logger:    logger,
		Router:    r,
		Validator: validator,
		Ch:        ch,
	})

	err := http.ListenAndServe("0.0.0.0:3000", r)
	if err != nil {
		panic(err)
	}
}
