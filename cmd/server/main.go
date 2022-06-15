package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/users"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/validator"
	"github.com/rs/zerolog/log"
)

type config struct {
	port int
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "Server port")

	usersService := &users.UsersService{
		Validate: validator.InitValidator(),
	}

	usersHandler := &users.UsersHandler{
		UsersService: usersService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandler.RegisterUser)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: mux,
	}

	log.Info().Msgf("Starting server on %s", server.Addr)
	log.Fatal().Err(server.ListenAndServe()).Msg("")
}
