package users

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func TestGivenAllFieldsAreSetWhenUpdateUserShouldReturnUpdatedUser(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	updateUserRequestData := UpdateUserRequest{}

	err = faker.FakeData(&updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(updatedUser.User.Email, *updateUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.User.Username != *updateUserRequestData.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, *updateUserRequestData.User.Username)
	}

	if updatedUser.User.Email != *updateUserRequestData.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, *updateUserRequestData.User.Email)
	}

	if updatedUser.User.Bio != *updateUserRequestData.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, *updateUserRequestData.User.Bio)
	}

	if updatedUser.User.Image != *updateUserRequestData.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, *updateUserRequestData.User.Image)
	}

	if updatedUser.User.Username != loggedUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, loggedUser.User.Username)
	}

	if updatedUser.User.Email != loggedUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, loggedUser.User.Email)
	}

	if updatedUser.User.Bio != loggedUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, loggedUser.User.Bio)
	}

	if updatedUser.User.Image != loggedUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, loggedUser.User.Image)
	}
}

func TestGivenUsernameIsSetWhenUpdateUserShouldReturnUpdatedUser(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	username := faker.Username()

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Username: &username,
		},
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(updatedUser.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.User.Username != *updateUserRequestData.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, *updateUserRequestData.User.Username)
	}

	if updatedUser.User.Email != registeredUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, registeredUser.User.Email)
	}

	if updatedUser.User.Bio != registeredUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, registeredUser.User.Bio)
	}

	if updatedUser.User.Image != registeredUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, registeredUser.User.Image)
	}

	if updatedUser.User.Username != loggedUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, loggedUser.User.Username)
	}

	if updatedUser.User.Email != loggedUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, loggedUser.User.Email)
	}

	if updatedUser.User.Bio != loggedUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, loggedUser.User.Bio)
	}

	if updatedUser.User.Image != loggedUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, loggedUser.User.Image)
	}
}

func TestGivenEmailIsSetWhenUpdateUserShouldReturnUpdatedUser(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	email := faker.Email()

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Email: &email,
		},
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(updatedUser.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.User.Username != registeredUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, registeredUser.User.Username)
	}

	if updatedUser.User.Email != *updateUserRequestData.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, *updateUserRequestData.User.Email)
	}

	if updatedUser.User.Bio != registeredUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, registeredUser.User.Bio)
	}

	if updatedUser.User.Image != registeredUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, registeredUser.User.Image)
	}

	if updatedUser.User.Username != loggedUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, loggedUser.User.Username)
	}

	if updatedUser.User.Email != loggedUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, loggedUser.User.Email)
	}

	if updatedUser.User.Bio != loggedUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, loggedUser.User.Bio)
	}

	if updatedUser.User.Image != loggedUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, loggedUser.User.Image)
	}
}

func TestGivenBioIsSetWhenUpdateUserShouldReturnUpdatedUser(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	bio := faker.Sentence()

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Bio: &bio,
		},
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(updatedUser.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.User.Username != registeredUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, registeredUser.User.Username)
	}

	if updatedUser.User.Email != registeredUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, registeredUser.User.Email)
	}

	if updatedUser.User.Bio != *updateUserRequestData.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, *updateUserRequestData.User.Bio)
	}

	if updatedUser.User.Image != registeredUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, registeredUser.User.Image)
	}

	if updatedUser.User.Username != loggedUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, loggedUser.User.Username)
	}

	if updatedUser.User.Email != loggedUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, loggedUser.User.Email)
	}

	if updatedUser.User.Bio != loggedUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, loggedUser.User.Bio)
	}

	if updatedUser.User.Image != loggedUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, loggedUser.User.Image)
	}
}

func TestGivenImageIsSetWhenUpdateUserShouldReturnUpdatedUser(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	image := faker.URL()

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Image: &image,
		},
	}

	updatedUser, err := UpdateUserAndDecode(registeredUser.User.Token, updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	loggedUser, err := LoginAndDecode(updatedUser.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	if updatedUser.User.Username != registeredUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, registeredUser.User.Username)
	}

	if updatedUser.User.Email != registeredUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, registeredUser.User.Email)
	}

	if updatedUser.User.Bio != registeredUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, registeredUser.User.Bio)
	}

	if updatedUser.User.Image != *updateUserRequestData.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, *updateUserRequestData.User.Image)
	}

	if updatedUser.User.Username != loggedUser.User.Username {
		t.Fatalf("got %s, want %s", updatedUser.User.Username, loggedUser.User.Username)
	}

	if updatedUser.User.Email != loggedUser.User.Email {
		t.Fatalf("got %s, want %s", updatedUser.User.Email, loggedUser.User.Email)
	}

	if updatedUser.User.Bio != loggedUser.User.Bio {
		t.Fatalf("got %s, want %s", updatedUser.User.Bio, loggedUser.User.Bio)
	}

	if updatedUser.User.Image != loggedUser.User.Image {
		t.Fatalf("got %s, want %s", updatedUser.User.Image, loggedUser.User.Image)
	}
}

func TestGivenUserDoesNotExistsWhenUpdateUserShouldReturnNotFound(t *testing.T) {
	const nonExistentUserToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4MTE2Mjg2MjcsImlhdCI6MTY1NTg2ODYyNywic3ViIjoibmVpbHBlYXJ0In0.0aDiR6d8jJsn9Ii9T176GhF34CqVT-KgTrU77BBjIgM"

	updateUserRequestData := UpdateUserRequest{}

	err := faker.FakeData(&updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := UpdateUser(nonExistentUserToken, updateUserRequestData)

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}

func TestGivenTokenIsInvalidWhenUpdateUserShouldReturnUnauthorized(t *testing.T) {
	const invalidToken = " "

	updateUserRequestData := UpdateUserRequest{}

	err := faker.FakeData(&updateUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := UpdateUser(invalidToken, updateUserRequestData)

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestGivenUsernameAlreadyExistsWhenUpdateUserShouldReturnUnprocessableEntity(t *testing.T) {
	existingUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&existingUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	existingUser, err := RegisterUserAndDecode(existingUserRequestData.User.Username, existingUserRequestData.User.Email, existingUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	registerUserRequestData := RegisterUserRequest{}

	err = faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Username: &existingUser.User.Username,
		},
	}

	response, err := UpdateUser(registeredUser.User.Token, updateUserRequestData)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	defer response.Body.Close()

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "User already exists" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "User already exists")
	}
}

func TestEmailIsTakenWhenUpdateUserShouldReturnUnprocessableEntity(t *testing.T) {
	existingUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&existingUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	existingUser, err := RegisterUserAndDecode(existingUserRequestData.User.Username, existingUserRequestData.User.Email, existingUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	registerUserRequestData := RegisterUserRequest{}

	err = faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Email: &existingUser.User.Email,
		},
	}

	response, err := UpdateUser(registeredUser.User.Token, updateUserRequestData)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	defer response.Body.Close()

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Email is taken" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Email is taken")
	}
}

func TestGivenEmailIsInvalidWhenUpdateUserShouldReturnUnprocessableEntity(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	email := "invalid"

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Email: &email,
		},
	}

	response, err := UpdateUser(registeredUser.User.Token, updateUserRequestData)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	defer response.Body.Close()

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Invalid email" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Invalid email")
	}
}

func TestGivenPasswordShorterThanEightCharactersWhenUpdateUserShouldReturnUnprocessableEntity(t *testing.T) {
	registerUserRequestData := RegisterUserRequest{}

	err := faker.FakeData(&registerUserRequestData)
	if err != nil {
		t.Fatal(err)
	}

	registeredUser, err := RegisterUserAndDecode(registerUserRequestData.User.Username, registerUserRequestData.User.Email, registerUserRequestData.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	password := "1234567"

	updateUserRequestData := UpdateUserRequest{
		User: updateUserRequestUser{
			Password: &password,
		},
	}

	response, err := UpdateUser(registeredUser.User.Token, updateUserRequestData)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got %d, want %d", response.StatusCode, http.StatusUnprocessableEntity)
	}

	defer response.Body.Close()

	responseData := &ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	if responseData.Errors.Body[0] != "Password must contain at least 8 characters" {
		t.Fatalf("got %s, want %s", responseData.Errors.Body[0], "Password must contain at least 8 characters")
	}
}
