package dto

type UserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}
type UserGenerateToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
