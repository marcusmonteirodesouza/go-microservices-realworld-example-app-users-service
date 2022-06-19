package users

type User struct {
	Username     string
	Email        string
	PasswordHash string
	Bio          *string
	Image        *string
}

func NewUser(username, email string, passwordHash string, bio *string, image *string) User {
	return User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio,
		Image:        image,
	}
}
