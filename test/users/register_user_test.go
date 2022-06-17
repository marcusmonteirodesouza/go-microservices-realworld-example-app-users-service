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

const apiUrl = "http://localhost:8080/users"
const contentType = "application/json"

type RegisterUserRequest struct {
	User struct {
		Username string `json:"username" faker:"username"`
		Email    string `json:"email" faker:"email"`
		Password string `json:"password" faker:"password"`
	} `json:"user"`
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

func TestGivenValidRequestShouldReturnUser(t *testing.T) {
	requestData := &RegisterUserRequest{}
	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Post(apiUrl, contentType, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("got %d, want %d", response.StatusCode, http.StatusCreated)
	}

	responseData := &RegisterUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.User.Username != requestData.User.Username {
		t.Errorf("got %s, want %s", responseData.User.Username, requestData.User.Username)
	}

	if responseData.User.Email != requestData.User.Email {
		t.Errorf("got %s, want %s", responseData.User.Email, requestData.User.Email)
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
		t.Errorf("got %s, want %s", claims.Subject, requestData.User.Username)
	}

	now := time.Now().Unix()
	if claims.ExpiresAt <= now {
		t.Errorf("ExpiresAt claim must be greater than now. got %d, now %d", claims.ExpiresAt, now)
	}

	if claims.IssuedAt-now > 60 {
		t.Errorf("IssuedAt claim must be close to now. got %d, now %d", claims.IssuedAt, now)
	}

	if claims.IssuedAt+int64(jwtSecondsToExpire) != claims.ExpiresAt {
		t.Errorf("got %d, want %d", claims.ExpiresAt, claims.IssuedAt+int64(jwtSecondsToExpire))
	}
}
