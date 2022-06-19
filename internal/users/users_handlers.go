package users

import (
	"encoding/json"
	"net/http"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/auth"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/errors"
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
	User *userResponseUser `json:"user"`
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
		User: &userResponseUser{
			Email:    email,
			Token:    token,
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
	if r.Method != "POST" {
		methodNotAllowed(w, r)
		return
	}

	var request struct {
		User struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		unprocessableEntity(w, r, []error{err})
		return
	}

	user, err := h.UsersService.RegisterUser(r.Context(), request.User.Username, request.User.Email, request.User.Password)
	if err != nil {
		if _, ok := err.(*errors.InvalidArgumentError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		if _, ok := err.(*errors.AlreadyExistsError); ok {
			unprocessableEntity(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *UsersHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		methodNotAllowed(w, r)
		return
	}

	var request struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		unprocessableEntity(w, r, []error{err})
		return
	}

	isValidPassword, err := h.UsersService.IsValidPassword(r.Context(), request.User.Email, request.User.Password)
	if err != nil || !isValidPassword {
		unauthorized(w, r)
		return
	}

	user, err := h.UsersService.GetUserByEmail(r.Context(), request.User.Email)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func (h *UsersHandlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		methodNotAllowed(w, r)
		return
	}

	username := r.Context().Value(auth.UsernameContextKey).(string)

	user, err := h.UsersService.GetUserByUsername(r.Context(), username)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			notFound(w, r, []error{err})
			return
		}

		internalServerError(w, r, err)
		return
	}

	token, err := h.JwtService.GenerateToken(user.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	responseBody := newUserResponse(user.Email, *token, user.Username, user.Bio, user.Image)

	response, err := json.Marshal(responseBody)
	if err != nil {
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

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
	log.Error().Err(err).Msg("")
	http.Error(w, "Internal Server Error", http.StatusUnauthorized)
}
