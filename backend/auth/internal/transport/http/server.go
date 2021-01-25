package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", endpoints.Auth.SignUp)
		authGroup.POST("/sign-in", endpoints.Auth.SignIn)
		authGroup.PUT("/token", endpoints.Auth.RefreshToken)
		authGroup.POST("/authenticate", endpoints.Auth.Authenticate)

		profileGroup := authGroup.Group("/user")
		{
			profileGroup.GET("/get-by-email", endpoints.Auth.SearchProfileByAnthroponym)
		}
	}

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
