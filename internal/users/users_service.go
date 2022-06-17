package users

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-playground/validator/v10"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersService struct {
	Validate  *validator.Validate
	Firestore *firestore.Client
}

const usersCollection = "users"

type userDocData struct {
	email         string
	password_hash string
	bio           string
	image         string
}

func (s *UsersService) RegisterUser(username string, email string, password string) (*User, error) {
	if len(strings.TrimSpace(username)) == 0 {
		return nil, &errors.InvalidArgumentError{Message: "Username cannot be blank"}
	}

	err := s.Validate.Var(email, "email")
	if err != nil {
		return nil, &errors.InvalidArgumentError{Message: "Invalid email"}
	}

	if len(password) < 8 {
		return nil, &errors.InvalidArgumentError{Message: "Password must contain at least 8 characters"}
	}

	ctx := context.Background()
	userDocPath := fmt.Sprintf("%s/%s", usersCollection, username)
	userDocRef := s.Firestore.Doc(userDocPath)
	_, err = userDocRef.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, err
		}
	} else {
		return nil, &errors.AlreadyExistsError{Message: "User already exists"}
	}

	user := &User{
		Username: username,
		Email:    email,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	_, err = userDocRef.Create(ctx, &userDocData{
		email:         user.Email,
		password_hash: string(passwordHash),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
