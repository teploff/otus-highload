package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			profileGroup.GET("/get-id-by-token", endpoints.Auth.GetUserIDByAccessToken)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
