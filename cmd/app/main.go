package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/auth"
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

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		log.Fatal().Err(err).Msg("Environment variable 'JWT_SECRET_KEY' must be set and not be empty")
	}

	jwtSecondsToExpire, err := strconv.Atoi(os.Getenv("JWT_SECONDS_TO_EXPIRE"))
	if err != nil {
		log.Fatal().Err(err).Msg("Environment variable 'JWT_SECONDS_TO_EXPIRE' must be set and not be empty")
	}

	jwtService := auth.NewJwtService(jwtSecretKey, jwtSecondsToExpire)

	firestoreClient, err := firestore.InitFirestore(ctx, firestoreProjectId)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing the Firestore client")
	}

	defer firestoreClient.Close()

	usersService := users.NewUsersService(*validator.InitValidator(), *firestoreClient)

	usersHandlers := users.NewUsersHandlers(usersService, jwtService)

	authMiddleware := auth.NewAuthMiddleware(jwtService)

	router := chi.NewRouter()
	router.Post("/users", usersHandlers.RegisterUser)
	router.Post("/users/login", usersHandlers.Login)
	router.Get("/user", authMiddleware.Authenticate(usersHandlers.GetCurrentUser))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Info().Msgf("Starting server on %s", server.Addr)
	log.Fatal().Err(server.ListenAndServe()).Msg("")
}
