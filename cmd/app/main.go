package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/firestore"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/users"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/validator"
	"github.com/rs/zerolog/log"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Environment variable 'PORT' must be set and set to an integer")
	}

	firestoreProjectId := os.Getenv("FIRESTORE_PROJECT_ID")
	if len(firestoreProjectId) == 0 {
		log.Fatal().Err(err).Msg("Environment variable 'FIRESTORE_PROJECT_ID' must be set and not be empty")
	}

	ctx := context.Background()

	firestoreClient, err := firestore.InitFirestore(ctx, firestoreProjectId)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing the Firestore client")
	}

	defer firestoreClient.Close()

	usersService := &users.UsersService{
		Validate:  validator.InitValidator(),
		Firestore: firestoreClient,
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		log.Fatal().Err(err).Msg("Environment variable 'JWT_SECRET_KEY' must be set and not be empty")
	}

	jwtSecondsToExpire, err := strconv.Atoi(os.Getenv("JWT_SECONDS_TO_EXPIRE"))
	if err != nil {
		log.Fatal().Err(err).Msg("Environment variable 'JWT_SECONDS_TO_EXPIRE' must be set and not be empty")
	}

	jwtService := &users.JwtService{
		SecretKey:       jwtSecretKey,
		SecondsToExpire: jwtSecondsToExpire,
	}

	usersHandlers := &users.UsersHandlers{
		UsersService: usersService,
		JwtService:   jwtService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandlers.RegisterUser)
	mux.HandleFunc("/users/login", usersHandlers.Login)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Info().Msgf("Starting server on %s", server.Addr)
	log.Fatal().Err(server.ListenAndServe()).Msg("")
}
