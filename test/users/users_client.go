package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterUserRequest struct {
	User registerUserRequestUser `json:"user"`
}

type registerUserRequestUser struct {
	Username string `json:"username" faker:"username"`
	Email    string `json:"email" faker:"email"`
	Password string `json:"password" faker:"password"`
}

func NewRegisterUserRequest(username string, email string, password string) RegisterUserRequest {
	return RegisterUserRequest{
		User: registerUserRequestUser{
			Username: username,
			Email:    email,
			Password: password,
		},
	}
}

type LoginRequest struct {
	User loginRequestUser `json:"user"`
}

type loginRequestUser struct {
	Email    string `json:"email" faker:"email"`
	Password string `json:"password" faker:"password"`
}

func NewLoginRequest(email string, password string) LoginRequest {
	return LoginRequest{
		User: loginRequestUser{
			Email:    email,
			Password: password,
		},
	}
}

type UpdateUserRequest struct {
	User updateUserRequestUser `json:"user"`
}

type updateUserRequestUser struct {
	Username *string `json:"username" faker:"username"`
	Email    *string `json:"email" faker:"email"`
	Password *string `json:"password" faker:"password"`
	Bio      *string `json:"bio" faker:"paragraph"`
	Image    *string `json:"image" faker:"url"`
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

type GetUserResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
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

	requestData := NewRegisterUserRequest(username, email, password)

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

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusCreated)
	}

	defer response.Body.Close()

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func Login(email string, password string) (*http.Response, error) {
	const url = "http://localhost:8080/users/login"

	requestData := NewLoginRequest(email, password)

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

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusOK)
	}

	defer response.Body.Close()

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func GetCurrentUser(tokenString string) (*http.Response, error) {
	client := &http.Client{}
	const url = "http://localhost:8080/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Forwarded-Authorization", fmt.Sprintf("Bearer %s", tokenString))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetCurrentUserAndDecode(tokenString string) (*UserResponse, error) {
	response, err := GetCurrentUser(tokenString)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusOK)
	}

	defer response.Body.Close()

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func GetUserByUsername(username string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8080/users/%s", username)

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetUserByUsernameAndDecode(username string) (*GetUserResponse, error) {
	response, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusOK)
	}

	defer response.Body.Close()

	responseData := &GetUserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func UpdateUser(tokenString string, request UpdateUserRequest) (*http.Response, error) {
	client := &http.Client{}
	const url = "http://localhost:8080/user"

	requestBytes, err := json.Marshal(request)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Forwarded-Authorization", fmt.Sprintf("Bearer %s", tokenString))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func UpdateUserAndDecode(tokenString string, request UpdateUserRequest) (*UserResponse, error) {
	response, err := UpdateUser(tokenString, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d, want %d", response.StatusCode, http.StatusOK)
	}

	defer response.Body.Close()

	responseData := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
