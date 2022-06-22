package users

type User struct {
	Id           string
	Username     string
	Email        string
	PasswordHash string
	Bio          *string
	Image        *string
}

func NewUser(id string, username string, email string, passwordHash string, bio *string, image *string) User {
	return User{
		Id:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio,
		Image:        image,
	}
}
