package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHTTPServer(addr string, endpoints *AuthEndpoints) *http.Server {
	router := gin.Default()

	//router.Group("/")
	//{
	//	router.GET("/questionnaire", endpoints)
	//}

	router.Group("/auth")
	{
		router.POST("/sign_up", endpoints.SignUp)
		router.POST("/sign_in", endpoints.SignIn)
		router.PUT("/token", endpoints.RefreshToken)
	}

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
