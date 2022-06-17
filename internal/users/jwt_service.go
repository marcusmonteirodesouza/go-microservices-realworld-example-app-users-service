package users

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtService struct {
	SecretKey       string
	SecondsToExpire int
}

func (s *JwtService) GenerateJwtToken(user User) (*string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Username,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Second * time.Duration(s.SecondsToExpire)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
