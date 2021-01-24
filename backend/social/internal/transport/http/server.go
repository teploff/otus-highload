package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	profileGroup := router.Group("/profile")
	{
		profileSearchGroup := profileGroup.Group("/search")
		{
			profileSearchGroup.GET("/anthroponym", endpoints.Profile.Search.GetByAnthroponym)
		}
	}

	socialGroup := router.Group("/social")
	{
		socialGroup.GET("/ws", endpoints.Ws.Connect)

		friendshipGroup := socialGroup.Group("/friendship")
		{
			friendshipGroup.POST("/create", endpoints.Social.CreateFriendship)
			friendshipGroup.POST("/confirm", endpoints.Social.ConfirmFriendship)
			friendshipGroup.POST("/reject", endpoints.Social.RejectFriendship)
			friendshipGroup.POST("/split-up", endpoints.Social.BreakFriendship)
		}

		socialGroup.GET("/friends", endpoints.Social.GetFriends)
		socialGroup.GET("/followers", endpoints.Social.GetFollowers)
		socialGroup.GET("/news", endpoints.Social.GetNews)
	}

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
