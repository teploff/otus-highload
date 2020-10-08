package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHTTPServer(addr string, endpoints *AuthEndpoints) *http.Server {
	router := gin.Default()

	group := router.Group("/auth")
	{
		group.POST("/sign_up", endpoints.SignUp)
		group.POST("/sign_in", endpoints.SignIn)
		group.PUT("/token", endpoints.RefreshToken)
	}

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
