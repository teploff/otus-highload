package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(addr string, endpoints *Endpoints) *http.Server {
	router := gin.Default()
	router.Use(ginmiddleware.Handler("", middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})))

	messengerGroup := router.Group("/messenger")
	{
		messengerGroup.GET("/ws", endpoints.Messenger.WS)
	}

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// just for demonstration RED metrics for prometheus
	router.GET("/internal-server", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{})
	})

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
