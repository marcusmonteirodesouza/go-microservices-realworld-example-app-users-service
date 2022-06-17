package users

import (
	"encoding/json"
	"net/http"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/errors"
	"github.com/rs/zerolog/log"
)

type UsersHandlers struct {
	UsersService *UsersService
	JwtService   *JwtService
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

	user, err := h.UsersService.RegisterUser(request.User.Username, request.User.Email, request.User.Password)
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

	token, err := h.JwtService.GenerateJwtToken(*user)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	responseBody := &userResponse{
		User: &userResponseUser{
			Email:    user.Email,
			Token:    *token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}

	response, err := json.Marshal(responseBody)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

type errorResponse struct {
	Errors *errorResponseErrors `json:"errors"`
}

type errorResponseErrors struct {
	Body []string `json:"body"`
}

func toErrorResponse(errors []error) errorResponse {
	var body []string
	for _, err := range errors {
		body = append(body, err.Error())
	}

	return errorResponse{
		Errors: &errorResponseErrors{
			Body: body,
		},
	}
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func unprocessableEntity(w http.ResponseWriter, r *http.Request, errors []error) {
	response, err := json.Marshal(toErrorResponse(errors))
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
	w.WriteHeader(http.StatusInternalServerError)
}
