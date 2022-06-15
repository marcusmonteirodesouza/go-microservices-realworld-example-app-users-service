package users

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/errors"
)

type UsersService struct {
	Validate *validator.Validate
}

func (s *UsersService) RegisterUser(username string, email string, password string) (*User, error) {
	if len(strings.TrimSpace(username)) == 0 {
		return nil, &errors.InvalidArgumentError{Message: "username cannot be blank"}
	}

	err := s.Validate.Var(email, "email")
	if err != nil {
		return nil, &errors.InvalidArgumentError{Message: "Invalid email"}
	}

	if len(password) < 8 {
		return nil, &errors.InvalidArgumentError{Message: "password must contain at least 8 characters"}
	}

	return &User{
		Id:       uuid.New(),
		Username: username,
		Email:    email,
		Bio:      nil,
		Image:    nil,
	}, nil
}
