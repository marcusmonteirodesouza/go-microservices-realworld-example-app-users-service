package users

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/go-cmp/cmp"
)

func TestGivenUserExistsWhenGetCurrentUserShouldReturnUser(t *testing.T) {
	requestData := RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)

	user, err := GetCurrentUserAndDecode(registeredUser.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(user, registeredUser) {
		t.Fatalf(cmp.Diff(user, registeredUser))
	}
}

func TestGivenUserDoesNotExistsWhenGetCurrentUserShouldReturnNotFound(t *testing.T) {
	const nonExistentUserToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4MTE2Mjg2MjcsImlhdCI6MTY1NTg2ODYyNywic3ViIjoibmVpbHBlYXJ0In0.0aDiR6d8jJsn9Ii9T176GhF34CqVT-KgTrU77BBjIgM"

	response, err := GetCurrentUser(nonExistentUserToken)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}

func TestGivenTokenIsInvalidWhenGetCurrentUserShouldReturnUnauthorized(t *testing.T) {
	const invalidToken = " "

	response, err := GetCurrentUser(invalidToken)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}
