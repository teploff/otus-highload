package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	router.Use(ginhttp.Middleware(opentracing.GlobalTracer()))

	socialGroup := router.Group("/social")
	{
		socialGroup.GET("/ws", endpoints.WS)

		profileGroup := socialGroup.Group("/profile")
		{
			profileGroup.GET("/search-by-anthroponym", endpoints.Profile.SearchByAnthroponym)
		}

		friendshipGroup := socialGroup.Group("/friendship")
		{
			friendshipGroup.POST("/create", endpoints.Friendship.Create)
			friendshipGroup.POST("/confirm", endpoints.Friendship.Confirm)
			friendshipGroup.POST("/reject", endpoints.Friendship.Reject)
			friendshipGroup.POST("/split-up", endpoints.Friendship.SplitUp)
			friendshipGroup.GET("/get-friends", endpoints.Friendship.GetFriends)
			friendshipGroup.GET("/get-followers", endpoints.Friendship.GetFollowers)
		}

		newsGroup := socialGroup.Group("/news")
		{
			newsGroup.GET("/get-news", endpoints.News.GetNews)
		}
	}

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
