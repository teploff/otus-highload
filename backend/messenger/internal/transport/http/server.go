package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()

	messengerGroup := router.Group("/messenger")
	{
		messengerGroup.GET("/ws", endpoints.Messenger.WS)
	}

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
