package user_service

type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" `
}

type SignInUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpUserDTO struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

type UpdateUserDTO struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	OldPassword      string `json:"old_password"`
	NewPassword      string `json:"new_password"`
	RepeatedPassword string `json:"repeated_password"`
}
