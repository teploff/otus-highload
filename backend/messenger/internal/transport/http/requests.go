package http

type AuthorizationHeader struct {
	AccessToken string `json:"access_token" binding:"required" header:"Authorization"`
}

type WSRequest struct {
	AccessToken string `json:"token"  binding:"required" form:"token"`
}
