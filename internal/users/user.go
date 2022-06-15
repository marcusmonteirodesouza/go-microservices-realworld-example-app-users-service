package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Username  string
	Email     string
	Bio       *string
	Image     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
