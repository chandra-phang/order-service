package v1

type AuthenticateDTO struct {
	Success bool    `json:"success"`
	Result  UserDTO `json:"result"`
}

type UserDTO struct {
	User UserDetailDTO `json:"user"`
}

type UserDetailDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
