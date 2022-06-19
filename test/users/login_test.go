package users

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt/v4"
)

func TestGivenValidRequestWhenLoginShouldReturnUser(t *testing.T) {
	requestData := &RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(registeredUser.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if loggedUser.User.Username != registeredUser.User.Username {
		t.Fatalf("got %s, want %s", loggedUser.User.Username, registeredUser.User.Username)
	}

	if loggedUser.User.Email != registeredUser.User.Email {
		t.Fatalf("got %s, want %s", loggedUser.User.Email, registeredUser.User.Email)
	}

	if loggedUser.User.Bio != registeredUser.User.Bio {
		t.Fatalf("got %s, want %s", loggedUser.User.Bio, registeredUser.User.Bio)
	}

	if loggedUser.User.Image != registeredUser.User.Image {
		t.Fatalf("got %s, want %s", loggedUser.User.Image, registeredUser.User.Image)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		t.Fatal("Environment variable 'JWT_SECRET_KEY' must be set and not be empty")
	}

	jwtSecondsToExpire, err := strconv.Atoi(os.Getenv("JWT_SECONDS_TO_EXPIRE"))
	if err != nil {
		t.Fatal("Environment variable 'JWT_SECONDS_TO_EXPIRE' must be set and not be empty")
	}

	parsedToken, err := jwt.ParseWithClaims(loggedUser.User.Token, jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(jwtSecretKey))
		return b, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	claims := parsedToken.Claims.(*jwt.StandardClaims)

	if claims.Subject != loggedUser.User.Username {
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

func TestGivenEmailNotFoundShouldReturnUnauthorized(t *testing.T) {
	requestData := &RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	_, err = RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	anotherUserRequestData := &RegisterUserRequest{}
	err = faker.FakeData(&anotherUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := Login(anotherUserRequestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(bodyBytes)
	if bodyString != "Unauthorized" {
		t.Fatalf("got %s, want %s", bodyString, "Unauthorized")
	}
}

func TestGivenPasswordIsIncorrectShouldReturnUnauthorized(t *testing.T) {
	requestData := &RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	anotherUserRequestData := &RegisterUserRequest{}
	err = faker.FakeData(&anotherUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := Login(registeredUser.User.Email, anotherUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(bodyBytes)
	if bodyString != "Unauthorized" {
		t.Fatalf("got %s, want %s", bodyString, "Unauthorized")
	}
}
