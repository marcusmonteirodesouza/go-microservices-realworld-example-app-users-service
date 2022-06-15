package main

import (
	"net/http"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/users"
	"github.com/rs/zerolog/log"
)

func main() {
	usersService := &users.UsersService{}
	usersHandler := &users.UsersHandler{
		UsersService: usersService,
	}

	http.HandleFunc("/users", usersHandler.RegisterUser)
	log.Fatal().Err(http.ListenAndServe(":8080", nil)).Msg("")
}
