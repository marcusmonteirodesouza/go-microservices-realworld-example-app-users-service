package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/auth"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/custom_errors"
	"github.com/rs/zerolog/log"
)

type UsersHandlers struct {
	UsersService UsersService
	JwtService   auth.JwtService
}

func NewUsersHandlers(usersService UsersService, jwtService auth.JwtService) UsersHandlers {
	return UsersHandlers{
		UsersService: usersService,
		JwtService:   jwtService,
	}
}

type userResponse struct {
	User userResponseUser `json:"user"`
}

type userResponseUser struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

func newUserResponse(email string, token string, username string, bio *string, image *string) userResponse {
	return userResponse{
		User: userResponseUser{
			Email:    email,
			Token:    token,
			Username: username,
			Bio:      bio,
			Image:    image,
		},
	}
}

type getUserResponse struct {
	User getUserResponseUser `json:"user"`
}

type getUserResponseUser struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

func newGetUserResponse(email string, username string, bio *string, image *string) getUserResponse {
	return getUserResponse{
		User: getUserResponseUser{
			Email:    email,
			Username: username,
			Bio:      bio,
			Image:    image,
		},
	}
}

type errorResponse struct {
	Errors errorResponseErrors `json:"errors"`
}

type errorResponseErrors struct {
	Body []string `json:"body"`
}

func newErrorResponse(errors []error) errorResponse {
	var body []string
	for _, err := range errors {
		body = append(body, err.Error())
	}

	return errorResponse{
		Errors: errorResponseErrors{
			Body: body,
		},
	}
}

func (h *UsersHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		User struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Error().Err(err).Msgf("Error decoding request")
		unprocessableEntity(w, r, []error{err})
		return
	}

	log.Info().Msgf("Registering User %s, email: %s...", request.User.Username, request.User.Password)
	user, err := h.UsersService.RegisterUser(r.Context(), request.User.Username, request.User.Email, request.User.Password)
	if err != nil {
		log.Error().Err(err).Msgf("Error registering User %s, email %s!", request.User.Username, request.User.Email)
		if _, ok := err.(*custom_errors.InvalidArgumentError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		if _, ok := err.(*custom_errors.AlreadyExistsError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		log.Error().Err(err).Msgf("Error generating Token for User %s, email %s", request.User.Username, request.User.Email)
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Error().Err(err).Msgf("Error marshalling response for User %s, email %s", request.User.Username, request.User.Email)
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *UsersHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Error().Err(err).Msgf("Error decoding request")
		unprocessableEntity(w, r, []error{err})
		return
	}

	IsCorrectPassword, err := h.UsersService.IsCorrectPassword(r.Context(), request.User.Email, request.User.Password)
	if err != nil || !IsCorrectPassword {
		log.Error().Err(err).Msgf("Error verifying password for email %s", request.User.Email)
		unauthorized(w, r)
		return
	}

	user, err := h.UsersService.GetUserByEmail(r.Context(), request.User.Email)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting User by email %s", request.User.Email)
		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		log.Error().Err(err).Msgf("Error generating token for email %s", request.User.Email)
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Error().Err(err).Msgf("Error marshalling response body for email %s", request.User.Email)
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func (h *UsersHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(auth.UsernameContextKey).(string)

	user, err := h.UsersService.GetUserByUsername(r.Context(), username)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting current User %s", username)
		if _, ok := err.(*custom_errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		log.Error().Err(err).Msgf("Error generating token for User %s", username)
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Error().Err(err).Msgf("Error marshalling response body for User %s", username)
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func (h *UsersHandlers) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := h.UsersService.GetUserByUsername(r.Context(), username)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting User %s", username)
		if _, ok := err.(*custom_errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	responseBody := newGetUserResponse(user.Email, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Error().Err(err).Msgf("Error marshalling response body for User %s", username)
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func (h *UsersHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(auth.UsernameContextKey).(string)

	user, err := h.UsersService.GetUserByUsername(r.Context(), username)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting User %s", username)

		if _, ok := err.(*custom_errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	var request struct {
		User struct {
			Email    *string `json:"email"`
			Username *string `json:"username"`
			Password *string `json:"password"`
			Image    *string `json:"image"`
			Bio      *string `json:"bio"`
		} `json:"user"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Error().Err(err).Msgf("Error decoding request")
		unprocessableEntity(w, r, []error{err})
		return
	}

	userUpdate := UserUpdate{
		Username: request.User.Username,
		Email:    request.User.Email,
		Password: request.User.Password,
		Bio:      request.User.Bio,
		Image:    request.User.Image,
	}

	user, err = h.UsersService.UpdateUserByUsername(r.Context(), username, userUpdate)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating User %s", username)

		if _, ok := err.(*custom_errors.InvalidArgumentError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		if _, ok := err.(*custom_errors.AlreadyExistsError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		log.Error().Err(err).Msgf("Error generating token for User %s", username)
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Error().Err(err).Msgf("Error marshalling response body for User %s", username)
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func notFound(w http.ResponseWriter, r *http.Request, errors []error) {
	response, err := json.Marshal(newErrorResponse(errors))
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(response)
}

func unprocessableEntity(w http.ResponseWriter, r *http.Request, errors []error) {
	response, err := json.Marshal(newErrorResponse(errors))
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(response)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
