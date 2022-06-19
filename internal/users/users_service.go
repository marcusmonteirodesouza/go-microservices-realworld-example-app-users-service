package users

import (
	"context"
	"fmt"
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
	Validate  validator.Validate
	Firestore firestore.Client
}

func NewUsersService(validate validator.Validate, firestore firestore.Client) UsersService {
	return UsersService{
		Validate:  validate,
		Firestore: firestore,
	}
}

const usersCollectionName = "users"

type userDocData struct {
	Email        string  `firestore:"email"`
	PasswordHash string  `firestore:"password_hash"`
	Bio          *string `firestore:"bio"`
	Image        *string `firestore:"image"`
}

func newUserDocData(email string, passwordHash string, bio *string, image *string) userDocData {
	return userDocData{
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio,
		Image:        image,
	}
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

	userDocRef := s.Firestore.Doc(fmt.Sprintf("%s/%s", usersCollectionName, username))
	_, err = userDocRef.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, err
		}
	} else {
		return nil, &errors.AlreadyExistsError{Message: "User already exists"}
	}

	existingUser, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); !ok {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, &errors.AlreadyExistsError{Message: "Email is taken"}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	userData := newUserDocData(email, string(passwordHash), nil, nil)

	_, err = userDocRef.Create(ctx, userData)

	if err != nil {
		return nil, err
	}

	user := NewUser(username, userData.Email, userData.PasswordHash, nil, nil)

	return &user, nil
}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	userDocRef := s.Firestore.Doc(fmt.Sprintf("%s/%s", usersCollectionName, username))
	userDocSnapshot, err := userDocRef.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, &errors.NotFoundError{Message: "User not found"}
		}
		return nil, err
	}

	userData := userDocData{}
	err = userDocSnapshot.DataTo(&userData)
	if err != nil {
		return nil, err
	}

	user := NewUser(userDocSnapshot.Ref.ID, userData.Email, userData.PasswordHash, userData.Bio, userData.Image)

	return &user, nil
}

func (s *UsersService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	usersCollection := s.Firestore.Collection(usersCollectionName)
	query := usersCollection.Where("email", "==", email).Limit(1)
	userDocs := query.Documents(ctx)
	defer userDocs.Stop()
	for {
		userDocSnapshot, err := userDocs.Next()
		if err == iterator.Done {
			return nil, &errors.NotFoundError{Message: "User not found"}
		} else if err != nil {
			return nil, err
		}

		userData := userDocData{}
		err = userDocSnapshot.DataTo(&userData)
		if err != nil {
			return nil, err
		}

		user := NewUser(userDocSnapshot.Ref.ID, userData.Email, userData.PasswordHash, userData.Bio, userData.Image)

		return &user, nil
	}
}

func (s *UsersService) IsValidPassword(ctx context.Context, email string, password string) (bool, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
