package users

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func TestGivenUserExistsWhenGetUserByUsernameShouldReturnUser(t *testing.T) {
	requestData := RegisterUserRequest{}

	err := faker.FakeData(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(requestData.User.Username, requestData.User.Email, requestData.User.Password)

	updateUserRequestData := UpdateUserRequest{}

	err = faker.FakeData(&updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	user, err := GetUserByUsernameAndDecode(updatedUser.User.Username)
	if err != nil {
		t.Fatal(err)
	}

	if user.User.Username != updatedUser.User.Username {
		t.Fatalf("got %s, want %s", user.User.Username, updatedUser.User.Username)
	}

	if user.User.Email != updatedUser.User.Email {
		t.Fatalf("got %s, want %s", user.User.Email, updatedUser.User.Email)
	}

	if user.User.Bio != updatedUser.User.Bio {
		t.Fatalf("got %s, want %s", user.User.Bio, updatedUser.User.Bio)
	}

	if user.User.Image != updatedUser.User.Image {
		t.Fatalf("got %s, want %s", user.User.Image, updatedUser.User.Image)
	}
}

func TestGivenUserDoesNotExistsWhenGetUserByUsernameShouldReturnNotFound(t *testing.T) {
	nonExistentUserToken := faker.Username()

	response, err := GetUserByUsername(nonExistentUserToken)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}
