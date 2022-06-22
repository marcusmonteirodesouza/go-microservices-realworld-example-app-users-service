package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtService struct {
	SecretKey       string
	SecondsToExpire int
}

func NewJwtService(secretKey string, secondsToExpire int) JwtService {
	return JwtService{
		SecretKey:       secretKey,
		SecondsToExpire: secondsToExpire,
	}
}

func (s *JwtService) GenerateToken(username string) (*string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Second * time.Duration(s.SecondsToExpire)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *JwtService) GetClaims(tokenString string) (*jwt.StandardClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})

	if claims, ok := parsedToken.Claims.(*jwt.StandardClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
