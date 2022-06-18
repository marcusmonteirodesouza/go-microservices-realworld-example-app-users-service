package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RegisterUserRequest struct {
	User registerUserResponseUser `json:"user"`
}

type registerUserResponseUser struct {
	Username string `json:"username" faker:"username"`
	Email    string `json:"email" faker:"email"`
	Password string `json:"password" faker:"password"`
}

type UserResponse struct {
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

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	return response, nil
}

func RegisterUserAndDecode(username string, email string, password string) (*UserResponse, error) {
	response, err := RegisterUser(username, email, password)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("got %d, want %d", response.StatusCode, http.StatusCreated))
	}

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func Login(email string, password string) (*http.Response, error) {
	const url = "http://localhost:8080/users/login"

	requestData := &RegisterUserRequest{
		User: registerUserResponseUser{
			Email:    email,
			Password: password,
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	return response, nil
}

func LoginAndDecode(email string, password string) (*UserResponse, error) {
	response, err := Login(email, password)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("got %d, want %d", response.StatusCode, http.StatusCreated))
	}

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
