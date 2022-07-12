package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-playground/validator/v10"
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-users-service/internal/custom_errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
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
	Username     string  `firestore:"username"`
	Email        string  `firestore:"email"`
	PasswordHash string  `firestore:"password_hash"`
	Bio          *string `firestore:"bio"`
	Image        *string `firestore:"image"`
}

func newUserDocData(username string, email string, passwordHash string, bio *string, image *string) userDocData {
	return userDocData{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio,
		Image:        image,
	}
}

func (s *UsersService) RegisterUser(ctx context.Context, username string, email string, password string) (*User, error) {
	if len(strings.TrimSpace(username)) == 0 {
		return nil, &custom_errors.InvalidArgumentError{Message: "Username cannot be blank"}
	}

	err := s.Validate.Var(email, "email")
	if err != nil {
		return nil, &custom_errors.InvalidArgumentError{Message: "Invalid email"}
	}

	err = validatePassword(password)
	if err != nil {
		return nil, &custom_errors.InvalidArgumentError{Message: err.Error()}
	}

	existingUser, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		if _, ok := err.(*custom_errors.NotFoundError); !ok {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, &custom_errors.AlreadyExistsError{Message: "User already exists"}
	}

	existingUser, err = s.GetUserByEmail(ctx, email)
	if err != nil {
		if _, ok := err.(*custom_errors.NotFoundError); !ok {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, &custom_errors.AlreadyExistsError{Message: "Email is taken"}
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	userDocRef := s.Firestore.Collection(usersCollectionName).NewDoc()
	userData := newUserDocData(username, email, *passwordHash, nil, nil)

	_, err = userDocRef.Create(ctx, userData)
	if err != nil {
		return nil, err
	}

	user := NewUser(userDocRef.ID, userData.Username, userData.Email, userData.PasswordHash, userData.Bio, userData.Image)

	return &user, nil
}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	usersCollection := s.Firestore.Collection(usersCollectionName)
	query := usersCollection.Where("username", "==", username).Limit(1)
	userDocs := query.Documents(ctx)
	defer userDocs.Stop()
	for {
		userDocSnapshot, err := userDocs.Next()
		if err == iterator.Done {
			return nil, &custom_errors.NotFoundError{Message: "User not found"}
		} else if err != nil {
			return nil, err
		}

		userData := userDocData{}
		err = userDocSnapshot.DataTo(&userData)
		if err != nil {
			return nil, err
		}

		user := NewUser(userDocSnapshot.Ref.ID, userData.Username, userData.Email, userData.PasswordHash, userData.Bio, userData.Image)

		return &user, nil
	}
}

func (s *UsersService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	usersCollection := s.Firestore.Collection(usersCollectionName)
	query := usersCollection.Where("email", "==", email).Limit(1)
	userDocs := query.Documents(ctx)
	defer userDocs.Stop()
	for {
		userDocSnapshot, err := userDocs.Next()
		if err == iterator.Done {
			return nil, &custom_errors.NotFoundError{Message: "User not found"}
		} else if err != nil {
			return nil, err
		}

		userData := userDocData{}
		err = userDocSnapshot.DataTo(&userData)
		if err != nil {
			return nil, err
		}

		user := NewUser(userDocSnapshot.Ref.ID, userData.Username, userData.Email, userData.PasswordHash, userData.Bio, userData.Image)

		return &user, nil
	}
}

type UserUpdate struct {
	Username *string
	Email    *string
	Password *string
	Bio      *string
	Image    *string
}

func (s *UsersService) UpdateUserByUsername(ctx context.Context, username string, userUpdate UserUpdate) (*User, error) {
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if userUpdate.Username != nil && *userUpdate.Username != user.Username {
		existingUser, err := s.GetUserByUsername(ctx, *userUpdate.Username)
		if err != nil {
			if _, ok := err.(*custom_errors.NotFoundError); !ok {
				return nil, err
			}
		}
		if existingUser != nil {
			return nil, &custom_errors.AlreadyExistsError{Message: "User already exists"}
		}
		user.Username = *userUpdate.Username
	}

	if userUpdate.Email != nil && *userUpdate.Email != user.Email {
		err := s.Validate.Var(userUpdate.Email, "email")
		if err != nil {
			return nil, &custom_errors.InvalidArgumentError{Message: "Invalid email"}
		}
		existingUser, err := s.GetUserByEmail(ctx, *userUpdate.Email)
		if err != nil {
			if _, ok := err.(*custom_errors.NotFoundError); !ok {
				return nil, err
			}
		}
		if existingUser != nil {
			return nil, &custom_errors.AlreadyExistsError{Message: "Email is taken"}
		}
		user.Email = *userUpdate.Email
	}

	if userUpdate.Password != nil {
		err := validatePassword(*userUpdate.Password)
		if err != nil {
			return nil, &custom_errors.InvalidArgumentError{Message: err.Error()}
		}
		passwordHash, err := hashPassword(*userUpdate.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = *passwordHash
	}

	if userUpdate.Bio != nil {
		user.Bio = userUpdate.Bio
	}

	if userUpdate.Image != nil {
		err = s.Validate.Var(userUpdate.Image, "url")
		if err != nil {
			return nil, &custom_errors.InvalidArgumentError{Message: err.Error()}
		}
		user.Image = userUpdate.Image
	}

	userDocRef := s.Firestore.Doc(fmt.Sprintf("%s/%s", usersCollectionName, user.Id))
	userDocData := newUserDocData(user.Username, user.Email, user.PasswordHash, user.Bio, user.Image)
	_, err = userDocRef.Set(ctx, userDocData)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UsersService) IsCorrectPassword(ctx context.Context, email string, password string) (bool, error) {
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

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must contain at least 8 characters")
	}

	return nil
}

func hashPassword(password string) (*string, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	passwordHash := string(passwordHashBytes)
	return &passwordHash, nil
}
