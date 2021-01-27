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
	router.Use(AuthenticateMiddleware(endpoints.cfg.Auth.Addr))

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

	socialGroup := router.Group("/social")
	{
		socialGroup.GET("/ws", endpoints.Social.WS)

		profileGroup := socialGroup.Group("/profile")
		{
			profileGroup.GET("/search-by-anthroponym", endpoints.Social.Profile.SearchByAnthroponym)
		}

		friendshipGroup := socialGroup.Group("/friendship")
		{
			friendshipGroup.POST("/create", endpoints.Social.Friendship.Create)
			friendshipGroup.POST("/confirm", endpoints.Social.Friendship.Confirm)
			friendshipGroup.POST("/reject", endpoints.Social.Friendship.Reject)
			friendshipGroup.POST("/split-up", endpoints.Social.Friendship.SplitUp)
			friendshipGroup.GET("/get-friends", endpoints.Social.Friendship.GetFriends)
			friendshipGroup.GET("/get-followers", endpoints.Social.Friendship.GetFollowers)
		}

		newsGroup := socialGroup.Group("/news")
		{
			newsGroup.GET("/get-news", endpoints.Social.News.GetNews)
		}
	}

	messengerGroup := router.Group("/messenger")
	{
		messengerGroup.GET("/ws", endpoints.Messenger.WS)
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
