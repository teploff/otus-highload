package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func AuthenticateMiddleware(authAddr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := req.Header{
			"Accept":        "application/json",
			"Authorization": findToken(c),
		}

		type Request struct {
			Resource string `json:"resource"`
		}

		body := Request{
			Resource: c.Request.URL.Path,
		}

		r, err := req.Post("http://"+authAddr+"/auth/authenticate", header, req.BodyJSON(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())

			return
		}

		var response authenticateResponse
		if err = r.ToJSON(&response); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)

			return
		}

		if !response.IsAuthenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Permission Denied")

			return
		}

		c.Next()
	}
}

func findToken(c *gin.Context) string {
	if token := c.Request.Header.Get("Authorization"); token != "" {
		return token
	}

	// For WS
	var request WSRequest
	if err := c.BindQuery(&request); err != nil {
		return ""
	}

	return request.AccessToken
}
