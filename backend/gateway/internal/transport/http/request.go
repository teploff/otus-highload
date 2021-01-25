package http

type WSRequest struct {
	AccessToken string `json:"token"  binding:"required" form:"token"`
}
