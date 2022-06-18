package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/joho/godotenv/autoload"
)

const contentType = "application/json"

type RegisterUserRequest struct {
	User registerUserResponseUser `json:"user"`
}

type registerUserResponseUser struct {
	Username string `json:"username" faker:"username"`
	Email    string `json:"email" faker:"email"`
	Password string `json:"password" faker:"password"`
}

type RegisterUserResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

type ErrorResponse struct {
	Errors *ErrorResponseErrors `json:"errors"`
}

type ErrorResponseErrors struct {
	Body []string `json:"body"`
}

func RegisterUser(username string, email string, password string) (*http.Response, error) {
	const url = "http://localhost:8080/users"

	requestData := &RegisterUserRequest{
		User: registerUserResponseUser{
			Username: username,
			Email:    email,
			Password: password,
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, contentType, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	return response, nil
}

func TestGivenValidRequestShouldReturnUser(t *testing.T) {
	requestData := &RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := RegisterUser(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusCreated)
	}

	responseData := &RegisterUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.User.Username != requestData.User.Username {
		t.Fatalf("got %s, want %s", responseData.User.Username, requestData.User.Username)
	}

	if responseData.User.Email != requestData.User.Email {
		t.Fatalf("got %s, want %s", responseData.User.Email, requestData.User.Email)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		t.Fatal("Environment variable 'JWT_SECRET_KEY' must be set and not be empty")
	}

	jwtSecondsToExpire, err := strconv.Atoi(os.Getenv("JWT_SECONDS_TO_EXPIRE"))
	if err != nil {
		t.Fatal("Environment variable 'JWT_SECONDS_TO_EXPIRE' must be set and not be empty")
	}

	parsedToken, err := jwt.ParseWithClaims(responseData.User.Token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(jwtSecretKey))
		return b, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	claims := parsedToken.Claims.(*jwt.StandardClaims)

	if claims.Subject != responseData.User.Username {
		t.Fatalf("got %s, want %s", claims.Subject, requestData.User.Username)
	}

	now := time.Now().Unix()
	if claims.ExpiresAt <= now {
		t.Fatalf("ExpiresAt claim must be greater than now. got %d, now %d", claims.ExpiresAt, now)
	}

	if claims.IssuedAt-now > 60 {
		t.Fatalf("IssuedAt claim must be close to now. got %d, now %d", claims.IssuedAt, now)
	}

	if claims.IssuedAt+int64(jwtSecondsToExpire) != claims.ExpiresAt {
		t.Fatalf("got %d, want %d", claims.ExpiresAt, claims.IssuedAt+int64(jwtSecondsToExpire))
	}
}

func TestGivenNoUsernameShouldReturnUnprocessableEntity(t *testing.T) {
	requestData := &RegisterUserRequest{}
	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}
	requestData.User.Username = " "

	response, err := RegisterUser(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Username cannot be blank" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Username cannot be blank")
	}
}

func TestGivenInvalidEmailShouldReturnUnprocessableEntity(t *testing.T) {
	requestData := &RegisterUserRequest{}
	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}
	requestData.User.Email = "invalid"

	response, err := RegisterUser(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Invalid email" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Invalid email")
	}
}

func TestGivenPasswordContainsLessThanEightCharactersShouldReturnUnprocessableEntity(t *testing.T) {
	requestData := &RegisterUserRequest{}
	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}
	requestData.User.Password = "1234567"

	response, err := RegisterUser(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Password must contain at least 8 characters" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Password must contain at least 8 characters")
	}
}

func TestGivenUsernameAlreadyExistsShouldReturnUnprocessableEntity(t *testing.T) {
	existingUserRequestData := &RegisterUserRequest{}
	err := faker.FakeData(&existingUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := RegisterUser(existingUserRequestData.User.Username, existingUserRequestData.User.Email, existingUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	anotherUserRequestData := &RegisterUserRequest{}
	err = faker.FakeData(&anotherUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err = RegisterUser(existingUserRequestData.User.Username, anotherUserRequestData.User.Email, anotherUserRequestData.User.Password)

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "User already exists" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "User already exists")
	}
}

func TestGivenEmailIsTakenShouldReturnUnprocessableEntity(t *testing.T) {
	existingUserRequestData := &RegisterUserRequest{}
	err := faker.FakeData(&existingUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := RegisterUser(existingUserRequestData.User.Username, existingUserRequestData.User.Email, existingUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	anotherUserRequestData := &RegisterUserRequest{}
	err = faker.FakeData(&anotherUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err = RegisterUser(anotherUserRequestData.User.Username, existingUserRequestData.User.Email, anotherUserRequestData.User.Password)

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Email is taken" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Email is taken")
	}
}
