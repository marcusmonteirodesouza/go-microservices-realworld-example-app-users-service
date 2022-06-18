package users

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-playground/validator/v10"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersService struct {
	Validate  *validator.Validate
	Firestore *firestore.Client
}

const usersCollectionName = "users"

type userDocData struct {
	Email        string `firestore:"email"`
	PasswordHash string `firestore:"password_hash"`
	Bio          string `firestore:"bio"`
	Image        string `firestore:"image"`
}

func (s *UsersService) RegisterUser(ctx context.Context, username string, email string, password string) (*User, error) {
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

	usersCollection := s.Firestore.Collection(usersCollectionName)
	userDocRef := usersCollection.Doc(username)
	_, err = userDocRef.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, err
		}
	} else {
		return nil, &errors.AlreadyExistsError{Message: "User already exists"}
	}

	existingUser, err := s.getUserByEmail(ctx, email)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); !ok {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, &errors.AlreadyExistsError{Message: "Email is taken"}
	}

	user := User{
		Username: username,
		Email:    email,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	_, err = userDocRef.Create(ctx, userDocData{
		Email:        user.Email,
		PasswordHash: string(passwordHash),
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UsersService) getUserByEmail(ctx context.Context, email string) (*User, error) {
	usersCollection := s.Firestore.Collection(usersCollectionName)
	query := usersCollection.Where("email", "==", email).Limit(1)
	docs := query.Documents(ctx)
	defer docs.Stop()
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			return nil, &errors.NotFoundError{Message: "User not found"}
		} else if err != nil {
			return nil, err
		}

		userData := userDocData{}
		err = doc.DataTo(&userData)
		if err != nil {
			return nil, err
		}

		return &User{
			Username: doc.Ref.ID,
			Email:    userData.Email,
			Bio:      &userData.Bio,
			Image:    &userData.Image,
		}, nil
	}
}
