package http

import (
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	router.Use(ginhttp.Middleware(opentracing.GlobalTracer()))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", endpoints.Auth.SignUp)
		authGroup.POST("/sign-in", endpoints.Auth.SignIn)
		authGroup.PUT("/token", endpoints.Auth.RefreshToken)
		authGroup.POST("/authenticate", endpoints.Auth.Authenticate)

		profileGroup := authGroup.Group("/user")
		{
			profileGroup.GET("/get-by-anthroponym", endpoints.Auth.GetUsersByAnthroponym)
			profileGroup.GET("/get-by-token", endpoints.Auth.GetUserByAccessToken)
			profileGroup.POST("/get-by-ids", endpoints.Auth.GetUserByIDs)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
