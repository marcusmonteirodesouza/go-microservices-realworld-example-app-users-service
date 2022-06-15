package users

import (
	"github.com/google/uuid"
)

type UsersService struct {
}

func (s *UsersService) RegisterUser(username string, email string, password string) (*User, error) {
	return &User{
		Id:       uuid.New(),
		Username: username,
		Email:    email,
		Bio:      nil,
		Image:    nil,
	}, nil
}
