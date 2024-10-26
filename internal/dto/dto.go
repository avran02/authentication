package dto

type RegisterRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email,omitempty"`
}

type RegisterResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type LogoutRequest struct {
	AccessToken string `json:"accessToken"`
}

type LogoutResponse struct {
	OK bool `json:"ok"`
}
