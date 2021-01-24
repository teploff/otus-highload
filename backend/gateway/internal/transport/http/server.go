package http

import (
	"github.com/swaggo/swag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization", "Access-Control-Expose-Headers")
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"PUT", "POST"}

	swag.Register(swag.Name, &swagDoc{})

	router.Use(cors.New(config))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", endpoints.Auth.SignUp)
		authGroup.POST("/sign-in", endpoints.Auth.SignIn)
		authGroup.PUT("/token", endpoints.Auth.RefreshToken)

		profileGroup := authGroup.Group("/user")
		{
			profileGroup.GET("/get-by-email", endpoints.Auth.GetUserIDByEmail)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

type swagDoc struct {
	doc string
}

func (s *swagDoc) ReadDoc() string {
	if s.doc == "" {
		data, err := ioutil.ReadFile("./api/swagger/swagger.json")
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		s.doc = string(data)
	}
	return s.doc
}
