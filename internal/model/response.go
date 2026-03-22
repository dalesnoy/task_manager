package model

// ErrorResponse — единый формат ошибок
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// TokenResponse — ответ при успешном логине
type TokenResponse struct {
	Token string `json:"token"`
}
