package http

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EmptyResponse struct {
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
