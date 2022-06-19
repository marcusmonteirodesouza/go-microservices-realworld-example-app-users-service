package users

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/go-cmp/cmp"
)

func TestGivenUserExistsWhenValidRequestShouldReturnUser(t *testing.T) {
	requestData := &RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)

	user, err := GetUserAndDecode(registeredUser.User.Token)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(user, registeredUser) {
		t.Fatalf(cmp.Diff(user, registeredUser))
	}
}

func TestGivenUserDoesNotExistsShouldReturnUnauthorized(t *testing.T) {
	const nonExistentUserToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU3NTM2NjQsImlhdCI6MTY1NTY2NzI2NCwic3ViIjoiMTc5ZGM1NzktMjFjZS00Mjg0LTg4YzItMjcyMjI3MzAyZjY0In0.3y6b232RRGvZYxgIoYwFb6l53KruHJhI392IbTRbA84"

	response, err := GetUser(nonExistentUserToken)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestGivenTokenIsInvalidShouldReturnUnauthorized(t *testing.T) {
	const nonExistentUserToken = " "

	response, err := GetUser(nonExistentUserToken)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}
