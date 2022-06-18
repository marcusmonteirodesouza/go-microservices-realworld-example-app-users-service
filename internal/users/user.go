package users

type User struct {
	Username     string
	Email        string
	PasswordHash string
	Bio          *string
	Image        *string
}
